package asg

import (
	"fmt"

	"github.com/MZCBBD/AWSScheduler/aws/common"
	"github.com/MZCBBD/AWSScheduler/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

type Asg struct {
	Id     string
	Status int
	Msg    string
	IsErr  bool
}

const (
	MsgStop         = "정상적으로 ASG가 중지되었습니다."
	MsgAlreadyStop  = "이미 ASG가 중지되어있습니다."
	MsgStart        = "정상적으로 ASG가 시작되었습니다."
	MsgAlreadyStart = "이미 ASG가 동작중입니다."
	MsgUnknown      = "알 수 없는 오류 입니다. ASG의 상태를 확인해주세요."
)

func New(id, msg string, status int, isErr bool) *Asg {
	return &Asg{Id: id, Msg: msg, Status: status, IsErr: isErr}
}

func (e *Asg) GetStatus() {
	sess := common.Sess()
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
}

func (e *Asg) Stop() {
	sess := common.Sess()
	svc := autoscaling.New(sess)

	if e.Status == 0 {
		e.Msg = MsgAlreadyStop
		e.IsErr = true
		fmt.Printf("Error: %t, CurrentStatus: %d, ID: %s, Msg: %s\n", e.IsErr, e.Status, e.Id, e.Msg)
	} else if e.Status >= 0 {
		input := &autoscaling.SetDesiredCapacityInput{
			AutoScalingGroupName: aws.String(e.Id),
			DesiredCapacity:      aws.Int64(0),
		}
		_, err := svc.SetDesiredCapacity(input)
		if err != nil {
			e.Msg = err.Error()
		}

		e.Msg = MsgStop
		e.Status = int(*input.DesiredCapacity)

		fmt.Printf("Error: %t, CurrentStatus: %d, ID: %s, Msg: %s\n", e.IsErr, e.Status, e.Id, e.Msg)
	} else {
		e.Msg = MsgUnknown
		e.IsErr = true
		fmt.Printf("Error: %t, CurrentStatus: %d, ID: %s, Msg: %s\n", e.IsErr, e.Status, e.Id, e.Msg)
	}

	utils.SlackNoti(e.Status, e.Id, e.Msg, e.IsErr)
}

func (e *Asg) Start() {
	sess := common.Sess()
	svc := autoscaling.New(sess)

	if e.Status >= 0 {
		e.Msg = MsgAlreadyStart
		e.IsErr = true
		fmt.Printf("Error: %t, CurrentStatus: %d, ID: %s, Msg: %s\n", e.IsErr, e.Status, e.Id, e.Msg)
	} else if e.Status == 0 {
		input := &autoscaling.SetDesiredCapacityInput{
			AutoScalingGroupName: aws.String(e.Id),
			DesiredCapacity:      aws.Int64(2),
		}

		_, err := svc.SetDesiredCapacity(input)
		if err != nil {
			e.Msg = err.Error()
		}

		e.Msg = MsgStart
		e.Status = int(*input.DesiredCapacity)

		fmt.Printf("Error: %t, CurrentStatus: %d, ID: %s, Msg: %s\n", e.IsErr, e.Status, e.Id, e.Msg)
	} else {
		e.Msg = MsgUnknown
		e.IsErr = true
		fmt.Printf("Error: %t, CurrentStatus: %d, ID: %s, Msg: %s\n", e.IsErr, e.Status, e.Id, e.Msg)
	}

	utils.SlackNoti(e.Status, e.Id, e.Msg, e.IsErr)
}
