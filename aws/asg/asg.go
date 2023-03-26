package asg

import (
	"fmt"

	"github.com/MZCBBD/AWSScheduler/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

type Asg struct {
	Id     string
	Status int
	Msg    string
}

const (
	StopMsg         = "정상적으로 ASG가 중지되었습니다."
	AlreadyStopMsg  = "이미 ASG가 중지되어있습니다."
	StartMsg        = "정상적으로 ASG가 시작되었습니다."
	AlreadyStartMsg = "이미 ASG가 동작중입니다."
	UnknownMsg      = "알수 없는 오류 입니다. ASG의 상태를 확인해주세요."
)

func New(id, msg string, status int) *Asg {
	return &Asg{Id: id, Msg: msg, Status: status}
}

func (e *Asg) GetStatus() error {
	sess := utils.Sess()
	svc := autoscaling.New(sess)

	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{aws.String(e.Id)},
	}
	output, err := svc.DescribeAutoScalingGroups(input)
	if err != nil {
		e.Msg = err.Error()
	}

	for _, autoScalingGroups := range output.AutoScalingGroups {
		e.Status = int(*autoScalingGroups.DesiredCapacity)
	}

	return nil
}

func (e *Asg) Stop() (string, error) {
	sess := utils.Sess()
	svc := autoscaling.New(sess)

	if e.Status == 0 {
		e.Msg = AlreadyStopMsg
		fmt.Printf("CurrentStatus: %d, ID: %s, Msg: %s\n", e.Status, e.Id, e.Msg)
	} else if e.Status >= 0 {
		input := &autoscaling.SetDesiredCapacityInput{
			AutoScalingGroupName: aws.String(e.Id),
			DesiredCapacity:      aws.Int64(0),
		}
		_, err := svc.SetDesiredCapacity(input)
		if err != nil {
			e.Msg = err.Error()
		}

		e.Msg = StopMsg
		e.Status = int(*input.DesiredCapacity)

		fmt.Printf("CurrentStatus: %d, ID: %s, Msg: %s\n", e.Status, e.Id, e.Msg)
		return "", nil
	} else {
		e.Msg = UnknownMsg
		fmt.Printf("CurrentStatus: %d, ID: %s, Msg: %s\n", e.Status, e.Id, e.Msg)
	}

	return "", nil
}

func (e *Asg) Start() (string, error) {
	sess := utils.Sess()
	svc := autoscaling.New(sess)

	if e.Status >= 0 {
		e.Msg = AlreadyStartMsg
		fmt.Printf("CurrentStatus: %d, ID: %s, Msg: %s\n", e.Status, e.Id, e.Msg)
	} else if e.Status == 0 {
		input := &autoscaling.SetDesiredCapacityInput{
			AutoScalingGroupName: aws.String(e.Id),
			DesiredCapacity:      aws.Int64(2),
		}

		_, err := svc.SetDesiredCapacity(input)
		if err != nil {
			e.Msg = err.Error()
		}

		e.Msg = StartMsg
		e.Status = int(*input.DesiredCapacity)

		fmt.Printf("CurrentStatus: %d, ID: %s, Msg: %s\n", e.Status, e.Id, e.Msg)
		return "", nil
	} else {
		e.Msg = UnknownMsg
		fmt.Printf("CurrentStatus: %d, ID: %s, Msg: %s\n", e.Status, e.Id, e.Msg)
	}

	return "", nil
}
