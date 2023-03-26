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
}

const (
	StopMsg         = "정상적으로 RDS 인스턴스가 중지되었습니다."
	AlreadyStopMsg  = "이미 RDS 인스턴스가 중지되어있습니다."
	StartMsg        = "정상적으로 RDS 인스턴스가 시작되었습니다."
	AlreadyStartMsg = "이미 RDS 인스턴스가 동작중입니다."
	UnknownMsg      = "알 수 없는 오류 입니다. 인스턴스의 상태를 확인해주세요."
)

func New(id, status, msg string) *Rds {
	return &Rds{Id: id, Status: status, Msg: msg}
}

func (e *Rds) GetStatus() error {
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

	return nil
}

func (e *Rds) Stop() (string, error) {
	sess := utils.Sess()
	svc := rds.New(sess)

	if e.Status == "stopped" {
		e.Msg = AlreadyStopMsg
		fmt.Printf("CurrentStatus: %s, ID: %s, Msg: %s\n", e.Status, e.Id, e.Msg)
	} else if e.Status == "available" {
		input := &rds.StopDBInstanceInput{
			DBInstanceIdentifier: &e.Id,
		}
		output, err := svc.StopDBInstance(input)
		utils.HandleErr(err)

		e.Msg = StopMsg
		e.Status = *output.DBInstance.DBInstanceStatus

		fmt.Printf("CurrentStatus: %s, ID: %s, Msg: %s\n", e.Status, e.Id, e.Msg)
		return "", nil
	} else {
		e.Msg = UnknownMsg
		fmt.Printf("CurrentStatus: %s, ID: %s, Msg: %s\n", e.Status, e.Id, e.Msg)
	}

	return "", nil
}

func (e *Rds) Start() (string, error) {
	sess := utils.Sess()
	svc := rds.New(sess)

	if e.Status == "available" {
		e.Msg = AlreadyStartMsg
		fmt.Printf("CurrentStatus: %s, ID: %s, Msg: %s\n", e.Status, e.Id, e.Msg)

		return "", nil
	} else if e.Status == "stopped" {
		input := &rds.StartDBInstanceInput{
			DBInstanceIdentifier: &e.Id,
		}
		output, err := svc.StartDBInstance(input)
		utils.HandleErr(err)

		e.Msg = StartMsg
		e.Status = *output.DBInstance.DBInstanceStatus

		fmt.Printf("CurrentStatus: %s, ID: %s, Msg: %s\n", e.Status, e.Id, e.Msg)
		return "", nil
	} else {
		e.Msg = UnknownMsg
		fmt.Printf("CurrentStatus: %s, ID: %s, Msg: %s\n", e.Status, e.Id, e.Msg)
	}

	return "", nil
}
