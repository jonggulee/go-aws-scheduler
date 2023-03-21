package main

import (
	"fmt"

	"github.com/MZCBBD/AWSScheduler/aws"
)

var services []string

func init() {
	services = append(services, "ec2", "rds")
}

func handler() {
	for _, service := range services {
		h := aws.NewAwsScheduler(service, "")
		fmt.Println(h.GetStatus())
	}
}

func main() {
	// lambda.Start(handler)

	// 로컬 테스트 실행 명령어
	handler()
}
