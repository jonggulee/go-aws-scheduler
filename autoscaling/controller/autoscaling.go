package autoscaling

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/slack-go/slack"
)

type autoScalingStatus struct {
	AutoScalingGroupName string
	DesiredCapacity      int
	Msg                  string
	Err                  error
}

func slackNoti(autoScalingStatus []autoScalingStatus) error {
	value := "```"
	successColor := "#18be52"
	failedColor := "#E96D76"
	failedValue := ":x: Failed"

	attachment := slack.Attachment{
		Fields: []slack.AttachmentField{
			{
				Title: "AutoScaling result",
				Value: ":white_check_mark: Success",
				Short: false,
			},
			{
				Title: "Target AutoScaling Groups",
				Value: "",
				Short: false,
			},
		},
	}
	attachment.Color = successColor

	for _, autoScalingStatus := range autoScalingStatus {

		if autoScalingStatus.Err != nil {
			attachment.Fields = append(attachment.Fields, slack.AttachmentField{
				Title: "Error",
				Value: fmt.Sprintf("```%s```", autoScalingStatus.Err),
				Short: false,
			})
			targetField := &attachment.Fields[0]
			targetField.Value = failedValue
			attachment.Color = failedColor
		}
		value += fmt.Sprintf("%s %s\n", autoScalingStatus.Msg, autoScalingStatus.AutoScalingGroupName)
	}
	attachment.Fields[1].Value = value + "```"

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

func getAutoScalingGroups(svc autoscalingiface.AutoScalingAPI, autoScalingGroupNames []*string) ([]autoScalingStatus, error) {
	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: autoScalingGroupNames,
	}
	describeAutoScalingGroupsOutput, err := svc.DescribeAutoScalingGroups(input)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	autoScalings := []autoScalingStatus{}
	for _, autoScalingGroup := range describeAutoScalingGroupsOutput.AutoScalingGroups {
		autoScaling := autoScalingStatus{
			AutoScalingGroupName: *autoScalingGroup.AutoScalingGroupName,
			DesiredCapacity:      int(*autoScalingGroup.DesiredCapacity),
		}
		autoScalings = append(autoScalings, autoScaling)
	}

	return autoScalings, nil
}

func StopAutoScaling(svc autoscalingiface.AutoScalingAPI, autoScalingGroupName string) error {
	input := &autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: aws.String(autoScalingGroupName),
		DesiredCapacity:      aws.Int64(0),
	}

	_, err := svc.SetDesiredCapacity(input)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func StopAutoScalingHandler() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := autoscaling.New(sess)

	autoScalingGroupNames := []*string{
		aws.String("eks-dev-mzwallet-eks-ng-2023011409042831340000000d-e2c2d7c2-f79e-e0b3-82a8-aa840609f926"),
		// aws.String("eks-dev-mzpay-eks-ng-20230313072707645300000002-26c36cee-d511-116b-a0b6-53db7f537fa1"),
	}

	autoScalingGroupStatus, err := getAutoScalingGroups(svc, autoScalingGroupNames)
	if err != nil {
		fmt.Println(err)
	}

	for _, autoScalingGroup := range autoScalingGroupStatus {
		if autoScalingGroup.DesiredCapacity >= 0 {
			err := StopAutoScaling(svc, autoScalingGroup.AutoScalingGroupName)
			if err != nil {
				fmt.Println(err)
			}
			slackNoti([]autoScalingStatus{autoScalingGroup})
		}
		if autoScalingGroup.DesiredCapacity >= 0 {
			fmt.Printf("이미 종료되었습니다. AutoScaling 상태를 확인해주세요.\n AutoScalingGroupName: %s / DesiredCapacity: %d\n", autoScalingGroup.AutoScalingGroupName, autoScalingGroup.DesiredCapacity)
		}
	}
}

func StartAutoScaling(svc autoscalingiface.AutoScalingAPI, autoScalingGroupName string) error {
	input := &autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: aws.String(autoScalingGroupName),
		DesiredCapacity:      aws.Int64(2),
	}

	_, err := svc.SetDesiredCapacity(input)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func StartAutoScalingHandler() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := autoscaling.New(sess)

	autoScalingGroupNames := []*string{
		aws.String("eks-dev-mzwallet-eks-ng-2023011409042831340000000d-e2c2d7c2-f79e-e0b3-82a8-aa840609f926"),
		// aws.String("eks-dev-mzpay-eks-ng-20230313072707645300000002-26c36cee-d511-116b-a0b6-53db7f537fa1"),
	}

	autoScalingGroupStatus, err := getAutoScalingGroups(svc, autoScalingGroupNames)
	if err != nil {
		fmt.Println(err)
	}

	for _, autoScalingGroup := range autoScalingGroupStatus {
		if autoScalingGroup.DesiredCapacity >= 0 {
			fmt.Printf("이미 동작 중입니다. AutoScaling 상태를 확인해주세요.\n AutoScalingGroupName: %s / DesiredCapacity: %d\n", autoScalingGroup.AutoScalingGroupName, autoScalingGroup.DesiredCapacity)
		}
		if autoScalingGroup.DesiredCapacity == 0 {
			err := StartAutoScaling(svc, autoScalingGroup.AutoScalingGroupName)
			if err != nil {
				fmt.Println(err)
			}
			slackNoti([]autoScalingStatus{autoScalingGroup})
		}
	}
}
