package autoscaling

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
)

type autoScalingStatus struct {
	AutoScalingGroupName string
	DesiredCapacity      int
	Msg                  string
	Err                  error
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
		aws.String("test"),
		aws.String("abcd"),
	}

	autoScalingStatus, err := getAutoScalingGroups(svc, autoScalingGroupNames)
	if err != nil {
		fmt.Println(err)
	}

	for _, autoScalingGroup := range autoScalingStatus {
		if autoScalingGroup.DesiredCapacity >= 0 {
			err := StopAutoScaling(svc, autoScalingGroup.AutoScalingGroupName)
			if err != nil {
				fmt.Println(err)
			}
		}
		if autoScalingGroup.DesiredCapacity >= 0 {
			fmt.Printf("이미 종료되었습니다. %s", autoScalingGroup.AutoScalingGroupName)
		}
	}

}
