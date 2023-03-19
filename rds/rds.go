package rds

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
	"github.com/slack-go/slack"
)

var AwsRegion = "ap-northeast-2"

type dbInstanceStatus struct {
	Status string
	Id     string
	Err    error
}

func slackNoti(inst []dbInstanceStatus) error {
	value := "```"
	color := "#18be52"
	attachment := slack.Attachment{
		Color: color,
		Fields: []slack.AttachmentField{
			{
				Title: "Stop DB result",
				Value: ":white_check_mark: Success",
				Short: false,
			},
			{
				Title: "Target",
				Value: "",
				Short: false,
			},
		},
	}
	for _, i := range inst {
		if i.Err != nil {
			newFieldValue := ":x: Failed"
			color = "#E96D76"
			attachment.Fields = append(attachment.Fields, slack.AttachmentField{
				Title: "Error",
				Value: fmt.Sprintf("```%s```", i.Err),
				Short: false,
			})
			targetField := &attachment.Fields[0]
			targetField.Value = newFieldValue
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

func getDBInstanceStatus(svc rdsiface.RDSAPI, dbInstanceIdentifier *string) (map[string]string, error) {
	instanceStatus := make(map[string]string)
	input := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: dbInstanceIdentifier,
	}

	output, err := svc.DescribeDBInstances(input)
	if err != nil {
		return instanceStatus, err
	}

	for _, dbInstance := range output.DBInstances {
		instanceStatus[*dbInstanceIdentifier] = string(*dbInstance.DBInstanceStatus)
	}
	return instanceStatus, nil
}

func stopDBInstance(svc rdsiface.RDSAPI, dbInstanceIdentifier *string) (string, error) {
	input := &rds.StopDBInstanceInput{
		DBInstanceIdentifier: dbInstanceIdentifier,
	}
	stopDBInstanceOutput, err := svc.StopDBInstance(input)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	stopDBInstanceStatus := stopDBInstanceOutput.DBInstance

	return *stopDBInstanceStatus.DBInstanceStatus, nil
}

func startDBInstance(svc rdsiface.RDSAPI, dbInstanceIdentifier *string) (string, error) {
	input := &rds.StartDBInstanceInput{
		DBInstanceIdentifier: dbInstanceIdentifier,
	}
	startDBInstanceOutput, err := svc.StartDBInstance(input)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	startDBInstanceStatus := startDBInstanceOutput.DBInstance

	return *startDBInstanceStatus.DBInstanceStatus, nil
}

func StopDBInstanceHandler() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := rds.New(sess)

	dbInstanceIdentifiers := []*string{
		aws.String("dev-mzwallet-mariadb"),
		aws.String("dev-mzpay-mariadb"),
	}

	instances := []dbInstanceStatus{}
	for _, dbInstanceIdentifier := range dbInstanceIdentifiers {
		getInstanceStatus, err := getDBInstanceStatus(svc, dbInstanceIdentifier)
		if err != nil {
			instance := dbInstanceStatus{
				Status: "RDS 정보를 찾을 수 없습니다. 종료하려는 RDS 이름을 확인해 주세요.",
				Id:     *dbInstanceIdentifier,
				Err:    err,
			}
			instances = append(instances, instance)
			continue
		}

		for InstanceId, status := range getInstanceStatus {
			if status == "stopped" {
				instance := dbInstanceStatus{
					Status: "이미 RDS가 중지되어 있습니다.",
					Id:     InstanceId,
					Err:    nil,
				}
				instances = append(instances, instance)
				continue
			}

			if status == "stopping" {
				instance := dbInstanceStatus{
					Status: "이미 RDS가 중지중 입니다. RDS 상태를 확인해주세요.",
					Id:     InstanceId,
					Err:    nil,
				}
				instances = append(instances, instance)
				continue
			}

			if status != "stopped" && status != "stopping" {
				stopDBInstanceStatusResult, err := stopDBInstance(svc, dbInstanceIdentifier)
				if err != nil {
					instance := dbInstanceStatus{
						Status: "RDS를 중지할 수 없습니다. RDS 정보를 확인해 주세요.",
						Id:     InstanceId,
						Err:    err,
					}
					instances = append(instances, instance)
					continue
				}
				if stopDBInstanceStatusResult == "stopping" {
					instance := dbInstanceStatus{
						Status: "정상적으로 RDS를 중지하였습니다.",
						Id:     InstanceId,
						Err:    err,
					}
					instances = append(instances, instance)
					// 확인 출력 코드
					fmt.Println("정상적으로 RDS를 중지하였습니다.", *dbInstanceIdentifier)
					continue
				}
			}
		}
	}
	slackNoti(instances)
}

func StartDBInstanceHandler() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := rds.New(sess)

	dbInstanceIdentifiers := []*string{
		aws.String("dev-mzwallet-mariadb"),
		aws.String("dev-mzpay-mariadb"),
	}

	instances := []dbInstanceStatus{}
	for _, dbInstanceIdentifier := range dbInstanceIdentifiers {
		getInstanceStatus, err := getDBInstanceStatus(svc, dbInstanceIdentifier)
		if err != nil {
			instance := dbInstanceStatus{
				Status: "RDS 정보를 찾을 수 없습니다. 종료하려는 RDS 이름을 확인해 주세요.",
				Id:     *dbInstanceIdentifier,
				Err:    err,
			}
			instances = append(instances, instance)
			continue
		}

		for InstanceId, status := range getInstanceStatus {
			if status == "available" {
				instance := dbInstanceStatus{
					Status: "이미 RDS가 동작하고 있습니다.",
					Id:     InstanceId,
					Err:    nil,
				}
				instances = append(instances, instance)
				continue
			}

			if status == "starting" {
				instance := dbInstanceStatus{
					Status: "이미 RDS가 시작중 입니다. RDS 상태를 확인해주세요.",
					Id:     InstanceId,
					Err:    nil,
				}
				instances = append(instances, instance)
				continue
			}

			if status != "available" && status != "starting" {
				startDBInstanceStatusResult, err := startDBInstance(svc, dbInstanceIdentifier)
				if err != nil {
					instance := dbInstanceStatus{
						Status: "RDS를 시작할 수 없습니다. RDS 정보를 확인해 주세요.",
						Id:     InstanceId,
						Err:    err,
					}
					instances = append(instances, instance)
					continue
				}
				if startDBInstanceStatusResult == "starting" {
					instance := dbInstanceStatus{
						Status: "정상적으로 RDS를 시작하였습니다.",
						Id:     InstanceId,
						Err:    nil,
					}
					instances = append(instances, instance)
					// 확인 출력 코드
					fmt.Println("정상적으로 RDS를 중지하였습니다.", *dbInstanceIdentifier)
					continue
				}
			}
		}
	}
	slackNoti(instances)
}
