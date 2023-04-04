package ec2

import (
	"fmt"

	"github.com/MZCBBD/AWSScheduler/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2 struct {
	Id     string
	Status string
	Msg    string
	IsErr  bool
}

const (
	MsgStop         = "정상적으로 인스턴스가 중지되었습니다."
	MsgAlreadyStop  = "이미 인스턴스가 중지되어있습니다."
	MsgStart        = "정상적으로 인스턴스가 시작되었습니다."
	MsgAlreadyStart = "이미 인스턴스가 동작중입니다."
	MsgUnknown      = "알 수 없는 오류 입니다. 인스턴스의 상태를 확인해주세요."
)

func New(id, status, msg string, isErr bool) *EC2 {
	return &EC2{Id: id, Status: status, Msg: msg, IsErr: isErr}
}

func (e *EC2) GetStatus() {
	sess := utils.Sess()
	svc := ec2.New(sess)

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(e.Id)},
	}
	output, err := svc.DescribeInstances(input)
	if err != nil {
		e.Msg = err.Error()
	}

	for _, reservation := range output.Reservations {
		for _, instance := range reservation.Instances {
			e.Status = aws.StringValue(instance.State.Name)
		}
	}
}

func (e *EC2) Stop() {
	sess := utils.Sess()
	svc := ec2.New(sess)

	if e.Status == "stopped" {
		e.Msg = MsgAlreadyStop
		e.IsErr = true
		fmt.Printf("Error: %t, CurrentStatus: %s, ID: %s, Msg: %s\n", e.IsErr, e.Status, e.Id, e.Msg)

	} else if e.Status == "running" {
		input := &ec2.StopInstancesInput{
			InstanceIds: []*string{aws.String(e.Id)},
			DryRun:      aws.Bool(true),
		}

		output, err := svc.StopInstances(input)
		awsErr, ok := err.(awserr.Error)

		if ok && awsErr.Code() == "DryRunOperation" {
			input.DryRun = aws.Bool(false)
			output, err = svc.StopInstances(input)
			utils.HandleErr(err)
		}

		for _, stoppingInstance := range output.StoppingInstances {
			currentState := stoppingInstance.CurrentState
			e.Msg = MsgStop
			e.Status = *currentState.Name

			fmt.Printf("Error: %t, CurrentStatus: %s, ID: %s, Msg: %s\n", e.IsErr, e.Status, e.Id, e.Msg)
		}
	} else {
		e.Msg = MsgUnknown
		e.IsErr = true
		fmt.Printf("Error: %t, CurrentStatus: %s, ID: %s, Msg: %s\n", e.IsErr, e.Status, e.Id, e.Msg)
	}

	utils.SlackNoti(e.Status, e.Id, e.Msg, e.IsErr)
}

func (e *EC2) Start() {
	sess := utils.Sess()
	svc := ec2.New(sess)

	if e.Status == "running" {
		e.Msg = MsgAlreadyStart
		e.IsErr = true
		fmt.Printf("Error: %t, CurrentStatus: %s, ID: %s, Msg: %s\n", e.IsErr, e.Status, e.Id, e.Msg)
	} else if e.Status == "stopped" {
		input := &ec2.StartInstancesInput{
			InstanceIds: []*string{aws.String(e.Id)},
			DryRun:      aws.Bool(true),
		}

		output, err := svc.StartInstances(input)
		awsErr, ok := err.(awserr.Error)

		if ok && awsErr.Code() == "DryRunOperation" {
			input.DryRun = aws.Bool(false)
			output, err = svc.StartInstances(input)
			utils.HandleErr(err)
		}

		for _, startingInstance := range output.StartingInstances {
			currentState := startingInstance.CurrentState
			e.Msg = MsgStart
			e.Status = *currentState.Name

			fmt.Printf("Error: %t, CurrentStatus: %s, ID: %s, Msg: %s\n", e.IsErr, e.Status, e.Id, e.Msg)
		}
	} else {
		e.Msg = MsgUnknown
		e.IsErr = true
		fmt.Printf("Error: %t, CurrentStatus: %s, ID: %s, Msg: %s\n", e.IsErr, e.Status, e.Id, e.Msg)
	}

	utils.SlackNoti(e.Status, e.Id, e.Msg, e.IsErr)
}
