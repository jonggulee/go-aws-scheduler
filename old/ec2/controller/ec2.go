package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/slack-go/slack"
)

var AwsRegion = "ap-northeast-2"

type instanceStatus struct {
	Status string
	Id     string
	Err    error
}

func slackNoti(inst []instanceStatus) error {
	value := "```"
	color := "#18be52"
	attachment := slack.Attachment{
		Fields: []slack.AttachmentField{
			{
				Title: "Instance start result",
				Value: ":white_check_mark: Success",
				Short: false,
			},
			{
				Title: "Target Instance",
				Value: "",
				Short: false,
			},
		},
	}
	for _, i := range inst {
		if i.Err != nil {
			newFieldValue := ":x: Failed"
			color = "#E96D76"
			attachment.Fields = append(attachment.Fields, slack.AttachmentField{
				Title: "Error",
				Value: fmt.Sprintf("```%s```", i.Err),
				Short: false,
			})
			targetField := &attachment.Fields[0]
			targetField.Value = newFieldValue
		}
		value += fmt.Sprintf("%s %s\n", i.Status, i.Id)
	}
	attachment.Fields[1].Value = value + "```"
	attachment.Color = color

	msg := slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment},
	}
	webhookUrl := "https://hooks.slack.com/services/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"
	err := slack.PostWebhook(webhookUrl, &msg)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func getInstanceStatus(svc ec2iface.EC2API, instanceIDs []*string) (map[string]string, error) {
	instanceStatus := make(map[string]string)
	input := &ec2.DescribeInstancesInput{
		InstanceIds: instanceIDs,
	}
	resp, err := svc.DescribeInstances(input)
	if err != nil {
		fmt.Println("인스턴스 정보 가져오기 실패:", err)
		return nil, err
	}

	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			state := instance.State
			instanceId := instance.InstanceId
			instanceStatus[*instanceId] = *state.Name
		}
	}

	return instanceStatus, nil
}

func startInstance(svc ec2iface.EC2API, instanceID []*string) (string, error) {
	input := &ec2.StartInstancesInput{
		InstanceIds: instanceID,
		DryRun:      aws.Bool(true),
	}
	startInstancesOutput, err := svc.StartInstances(input)

	awsErr, ok := err.(awserr.Error)

	if ok && awsErr.Code() == "DryRunOperation" {
		input.DryRun = aws.Bool(false)
		startInstancesOutput, err = svc.StartInstances(input)
		if err != nil {
			return "", err
		}
		for _, startingInstance := range startInstancesOutput.StartingInstances {
			currentState := startingInstance.CurrentState
			return *currentState.Name, nil
		}
	}
	for _, startingInstance := range startInstancesOutput.StartingInstances {
		currentState := startingInstance.CurrentState
		return *currentState.Name, nil
	}

	return "", err
}

func stopInstance(svc ec2iface.EC2API, instanceID []*string) (string, error) {
	input := &ec2.StopInstancesInput{
		InstanceIds: instanceID,
		DryRun:      aws.Bool(true),
	}

	stopInstancesOutput, err := svc.StopInstances(input)

	awsErr, ok := err.(awserr.Error)
	if ok && awsErr.Code() == "DryRunOperation" {
		input.DryRun = aws.Bool(false)
		stopInstancesOutput, err = svc.StopInstances(input)
		if err != nil {
			return "", err
		}
	}

	// EC2 인스턴스 종료 후 상태 확인
	for _, stoppingInstance := range stopInstancesOutput.StoppingInstances {
		currentState := stoppingInstance.CurrentState
		return *currentState.Name, nil
	}
	return "", err
}

func StartInstanceHandler() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := ec2.New(sess)

	instanceIDs := []*string{
		// sec-mgmt-jenkins
		aws.String("i-0acab82eb643706cd"),
		// sec-mgmt-eks-admin
		aws.String("i-04a4829ec4cf60254"),
	}

	getInstanceStatus, err := getInstanceStatus(svc, instanceIDs)
	if err != nil {
		fmt.Println("Check error the instance")
		fmt.Println(err)
	}

	instances := []instanceStatus{}
	for instanceId, status := range getInstanceStatus {
		if status == "running" {
			instnace := instanceStatus{
				Status: "이미 인스턴스가 실행중 입니다.",
				Id:     instanceId,
				Err:    err,
			}
			instances = append(instances, instnace)
			continue
		}
		if status == "pending" {
			instnace := instanceStatus{
				Status: "인스턴스가 pending 상태 입니다. 인스턴스의 상태를 확인해주세요.",
				Id:     instanceId,
				Err:    err,
			}
			instances = append(instances, instnace)
			continue
		}

		if status != "running" && status != "pending" {
			startInstanceResult, err := startInstance(svc, []*string{aws.String(instanceId)})
			if err != nil {
				instance := instanceStatus{
					Status: "인스턴스가 시작 중에 문제가 생겼습니다. 인스턴스의 상태를 확인해주세요.",
					Id:     instanceId,
					Err:    err,
				}
				instances = append(instances, instance)
				continue
			}
			if startInstanceResult == "pending" {
				instance := instanceStatus{
					Status: "정상적으로 인스턴스가 시작되었습니다.",
					Id:     instanceId,
					Err:    nil,
				}
				instances = append(instances, instance)
				continue
			}
		}
	}
	slackNoti(instances)
}

func StopInstanceHandler() {
	// 세션 생성
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// EC2 서비스 생성
	svc := ec2.New(sess)

	// EC2 인스턴스 ID 지정
	instanceIDs := []*string{
		// sec-mgmt-jenkins
		aws.String("i-0acab82eb643706cd"),
		// sec-mgmt-eks-admin
		aws.String("i-04a4829ec4cf60254"),
	}

	// EC2 인스턴스 상태 확인
	getInstanceStatus, err := getInstanceStatus(svc, instanceIDs)
	if err != nil {
		fmt.Println("Check error the instance")
		fmt.Println(err)
	}

	// 현재 EC2 인스턴스 상태에 따른 인스턴스 종료 여부 확인
	instances := []instanceStatus{}
	for instanceId, status := range getInstanceStatus {
		if status == "stopped" {
			instance := instanceStatus{
				Status: "이미 인스턴스가 중지되어있습니다.",
				Id:     instanceId,
				Err:    err,
			}
			instances = append(instances, instance)
			continue
		}

		if status == "stopping" {
			instance := instanceStatus{
				Status: "이미 인스턴스가 중지중 입니다. 인스턴스의 상태를 확인해주세요.",
				Id:     instanceId,
				Err:    err,
			}
			instances = append(instances, instance)
			continue
		}

		if status != "stopped" && status != "stopping" {
			stopInstanceResult, err := stopInstance(svc, []*string{aws.String(instanceId)})
			if err != nil {
				instance := instanceStatus{
					Status: "인스턴스가 중지 중에 문제가 생겼습니다. 인스턴스의 상태를 확인해주세요.",
					Id:     instanceId,
					Err:    err,
				}
				instances = append(instances, instance)
				continue
			}
			if stopInstanceResult == "stopping" {
				instance := instanceStatus{
					Status: "정상적으로 인스턴스가 중지되었습니다.",
					Id:     instanceId,
					Err:    nil,
				}
				instances = append(instances, instance)
				continue
			}
		}
	}
	slackNoti(instances)
}
