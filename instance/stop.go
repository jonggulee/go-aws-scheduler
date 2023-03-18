package instance

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/slack-go/slack"
)

type instanceStatus struct {
	Status string
	Id     string
	Err    error
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

func slackNoti(inst []instanceStatus) error {

	attachment := slack.Attachment{
		Color: "#18be52",
		Fields: []slack.AttachmentField{
			{
				Title: "Instance stop result",
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
	value := "```"
	color := "#18be52"
	for _, i := range inst {
		if i.Err != nil {
			color = "#E96D76"
			attachment.Fields = append(attachment.Fields, slack.AttachmentField{
				Title: "Error",
				Value: fmt.Sprintf("```%s```", i.Err),
				Short: false,
			})
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

func StopInstances() error {
	// 세션 생성
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// EC2 서비스 생성
	svc := ec2.New(sess)

	// EC2 인스턴스 ID 지정
	instanceIDs := []*string{
		aws.String("i-0420a345119580384"),
		aws.String("i-03ecc4d66213d0b37"),
	}

	// EC2 인스턴스 상태 확인
	getInstanceStatus, err := getInstanceStatus(svc, instanceIDs)
	if err != nil {
		fmt.Println("Check error the instance")
		fmt.Println(err)
		return err
	}

	// 현재 EC2 인스턴스 상태에 따른 인스턴스 종료 여부 확인
	instances := []instanceStatus{}
	for instanceId, status := range getInstanceStatus {
		if status == "stopped" {
			instance := instanceStatus{
				Status: "이미 인스턴스가 중지되어있습니다.",
				Id:     instanceId,
				Err:    nil,
			}
			instances = append(instances, instance)
			continue
		}

		if status == "Stopping" {
			instance := instanceStatus{
				Status: "이미 인스턴스가 중지중 입니다. 인스턴스의 상태를 확인해주세요.",
				Id:     instanceId,
				Err:    nil,
			}
			instances = append(instances, instance)
			continue
		}

		if status != "stopped" {
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
	return nil
}
