package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/MZCBBD/AWSScheduler/aws"
	"github.com/MZCBBD/AWSScheduler/utils"
)

func parseService(service string) map[string][]string {
	m := make(map[string][]string)
	for _, s := range strings.Split(service, ",") {
		parts := strings.Split(s, ":")
		key := parts[0]
		value := parts[1]
		if m[key] == nil {
			m[key] = []string{value}
		} else {
			m[key] = append(m[key], value)
		}
	}
	return m
}

func handler() {
	service := os.Getenv("service")
	action := os.Getenv("action")

	m := parseService(service)

	for service, IDs := range m {
		for _, Id := range IDs {
			scheduler := aws.NewAwsScheduler(service, Id)
			scheduler.GetStatus()
			if action == "stop" {
				fmt.Printf("Action: %s, Service: %s, ID: %s\n", action, service, Id)
				scheduler.Stop()
			}
			if action == "start" {
				fmt.Printf("Action: %s, Service: %s, ID: %s\n", action, service, Id)
				scheduler.Start()
			}
		}
	}
	utils.SendSlackMessage()
}

func main() {
	// lambda.Start(handler)

	// 로컬 테스트 실행 명령어
	handler()
}
