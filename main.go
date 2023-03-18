package main

import (
	"os"

	"github.com/MZCBBD/AWSScheduler/handlers"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler() {
	if os.Getenv("target") == "rds" {
		if os.Getenv("env") == "stop" {
			handlers.StopDBInstanceHandler()
		}
		// if os.Getenv("ENV") == "start" {
		// startDBInstanceHandler()
		// }
	}
}

func main() {
	// handler()
	lambda.Start(handler)
}
