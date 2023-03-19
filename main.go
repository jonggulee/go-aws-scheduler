package main

import (
	"os"

	ec2 "github.com/MZCBBD/AWSScheduler/ec2/controller"
	rds "github.com/MZCBBD/AWSScheduler/rds/controller"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler() {
	if os.Getenv("target") == "rds" {
		if os.Getenv("env") == "stop" {
			rds.StopDBInstanceHandler()
		}
		if os.Getenv("env") == "start" {
			rds.StartDBInstanceHandler()
		}
	}
	if os.Getenv("target") == "ec2" {
		if os.Getenv("env") == "stop" {
			ec2.StopInstanceHandler()
		}
		if os.Getenv("env") == "start" {
			ec2.StartInstanceHandler()
		}
	}
}

func main() {
	lambda.Start(handler)

	// 로컬 테스트 실행 명령어
	// handler()
}
