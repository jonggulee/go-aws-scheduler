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
}

const (
	StopMsg         = "정상적으로 인스턴스가 중지되었습니다."
	AlreadyStopMsg  = "이미 인스턴스가 중지되어있습니다."
	StartMsg        = "정상적으로 인스턴스가 시작되었습니다."
	AlreadyStartMsg = "이미 인스턴스가 동작중입니다."
	UnknownMsg      = "알수 없는 오류 입니다. 인스턴스의 상태를 확인해주세요."
)

func New(id, status, msg string) *EC2 {
	return &EC2{Id: id, Status: status, Msg: msg}
}

func (e *EC2) GetStatus() (string, error) {
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

	return e.Status, nil
}

func (e *EC2) Stop() (string, error) {
	sess := utils.Sess()
	svc := ec2.New(sess)

	if e.Status == "stopped" {
		e.Msg = AlreadyStopMsg
		fmt.Printf("CurrentStatus: %s, ID: %s, Msg: %s\n", e.Status, e.Id, e.Msg)

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
			e.Msg = StopMsg
			e.Status = *currentState.Name

			fmt.Printf("CurrentStatus: %s, ID: %s, Msg: %s\n", e.Status, e.Id, e.Msg)
			return "", nil
		}
	} else {
		e.Msg = UnknownMsg
		fmt.Printf("CurrentStatus: %s, ID: %s, Msg: %s\n", e.Status, e.Id, e.Msg)
	}

	return "", nil
}

func (e *EC2) Start() (string, error) {
	sess := utils.Sess()
	svc := ec2.New(sess)

	if e.Status == "running" {
		e.Msg = AlreadyStartMsg

		fmt.Printf("CurrentStatus: %s, ID: %s, Msg: %s\n", e.Status, e.Id, e.Msg)
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
			e.Msg = StartMsg
			e.Status = *currentState.Name

			fmt.Printf("CurrentStatus: %s, ID: %s, Msg: %s\n", e.Status, e.Id, e.Msg)
			return "", nil
		}
	} else {
		e.Msg = UnknownMsg
		fmt.Printf("CurrentStatus: %s, ID: %s, Msg: %s\n", e.Status, e.Id, e.Msg)
	}

	return "", nil
}
