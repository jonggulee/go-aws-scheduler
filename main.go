package main

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

func slackNoti(inst []instanceStatus) error {

	attachment := slack.Attachment{
		Color: "#18be52",
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
	// snippet-start:[ec2.go.start_stop_instances.start]
	input := &ec2.StartInstancesInput{
		InstanceIds: instanceID,
		DryRun:      aws.Bool(true),
	}
	startInstancesOutput, err := svc.StartInstances(input)

	awsErr, ok := err.(awserr.Error)

	if ok && awsErr.Code() == "DryRunOperation" {
		// Set DryRun to be false to enable starting the instances
		input.DryRun = aws.Bool(false)
		startInstancesOutput, err = svc.StartInstances(input)
		// snippet-end:[ec2.go.start_stop_instances.start]
		if err != nil {
			return "", err
		}

		return "", nil
	}
	for _, startingInstance := range startInstancesOutput.StartingInstances {
		currentState := startingInstance.CurrentState
		return *currentState.Name, nil
	}

	return "", err
}

func main() {
	// lambda.Start(myApp)
	// instance.MyApp()
	// lambda.Start(instance.MyApp())
	// instance.StopInstances()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := ec2.New(sess)

	instanceIDs := []*string{
		aws.String("i-0420a345119580384"),
		aws.String("i-03ecc4d66213d0b37"),
	}

	getInstanceStatus, err := getInstanceStatus(svc, instanceIDs)
	if err != nil {
		fmt.Println("Check error the instance")
		fmt.Println(err)
		// return err
	}

	instances := []instanceStatus{}
	for instanceId, status := range getInstanceStatus {
		if status == "running" {
			instnace := instanceStatus{
				Status: "이미 인스턴스가 실행중 입니다.",
				Id:     instanceId,
				Err:    nil,
			}
			instances = append(instances, instnace)
			continue
		}
		if status == "pending" {
			instnace := instanceStatus{
				Status: "인스턴스가 pending 상태 입니다. 인스턴스의 상태를 확인해주세요.",
				Id:     instanceId,
				Err:    nil,
			}
			instances = append(instances, instnace)
			continue
		}

		if status == "stopped" {
			startInstanceResult, err := startInstance(svc, []*string{aws.String(instanceId)})
			// err := startInstance(svc, instanceIDs)
			fmt.Println(startInstanceResult)
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

	// fmt.Println(instances)
	slackNoti(instances)

	// instance.StopInstances()
}
