package rds

import (
	"fmt"

	"github.com/MZCBBD/AWSScheduler/utils"
	"github.com/aws/aws-sdk-go/service/rds"
)

type Rds struct {
	Id     string
	Status string
	Msg    string
	IsErr  bool
}

const (
	MsgStop         = "정상적으로 RDS 인스턴스가 중지되었습니다."
	MsgAlreadyStop  = "이미 RDS 인스턴스가 중지되어있습니다."
	MsgStart        = "정상적으로 RDS 인스턴스가 시작되었습니다."
	MsgAlreadyStart = "이미 RDS 인스턴스가 동작중입니다."
	MsgUnknown      = "알 수 없는 오류 입니다. 인스턴스의 상태를 확인해주세요."
)

func New(id, status, msg string, isErr bool) *Rds {
	return &Rds{Id: id, Status: status, Msg: msg, IsErr: isErr}
}

func (e *Rds) GetStatus() {
	sess := utils.Sess()
	svc := rds.New(sess)

	input := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: &e.Id,
	}
	output, err := svc.DescribeDBInstances(input)
	if err != nil {
		e.Msg = err.Error()
	}

	for _, dbInstance := range output.DBInstances {
		e.Status = *dbInstance.DBInstanceStatus
	}
}

func (e *Rds) Stop() {
	sess := utils.Sess()
	svc := rds.New(sess)

	if e.Status == "stopped" {
		e.Msg = MsgAlreadyStop
		e.IsErr = true
		fmt.Printf("Error: %t, CurrentStatus: %s, ID: %s, Msg: %s\n", e.IsErr, e.Status, e.Id, e.Msg)

	} else if e.Status == "available" {
		input := &rds.StopDBInstanceInput{
			DBInstanceIdentifier: &e.Id,
		}
		output, err := svc.StopDBInstance(input)
		utils.HandleErr(err)

		e.Msg = MsgStop
		e.Status = *output.DBInstance.DBInstanceStatus

		fmt.Printf("Error: %t, CurrentStatus: %s, ID: %s, Msg: %s\n", e.IsErr, e.Status, e.Id, e.Msg)
	} else {
		e.Msg = MsgUnknown
		e.IsErr = true
		fmt.Printf("Error: %t, CurrentStatus: %s, ID: %s, Msg: %s\n", e.IsErr, e.Status, e.Id, e.Msg)
	}

	utils.SlackNoti(e.Status, e.Id, e.Msg, e.IsErr)
}

func (e *Rds) Start() {
	sess := utils.Sess()
	svc := rds.New(sess)

	if e.Status == "available" {
		e.Msg = MsgAlreadyStart
		e.IsErr = true

		fmt.Printf("Error: %t, CurrentStatus: %s, ID: %s, Msg: %s\n", e.IsErr, e.Status, e.Id, e.Msg)
	} else if e.Status == "stopped" {
		input := &rds.StartDBInstanceInput{
			DBInstanceIdentifier: &e.Id,
		}

		output, err := svc.StartDBInstance(input)
		utils.HandleErr(err)

		e.Msg = MsgStart
		e.Status = *output.DBInstance.DBInstanceStatus

		fmt.Printf("Error: %t, CurrentStatus: %s, ID: %s, Msg: %s\n", e.IsErr, e.Status, e.Id, e.Msg)
	} else {
		e.Msg = MsgUnknown
		e.IsErr = true
		fmt.Printf("Error: %t, CurrentStatus: %s, ID: %s, Msg: %s\n", e.IsErr, e.Status, e.Id, e.Msg)
	}

	utils.SlackNoti(e.Status, e.Id, e.Msg, e.IsErr)
}
