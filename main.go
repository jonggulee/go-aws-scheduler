package main

import (
	"os"

	rds "github.com/MZCBBD/AWSScheduler/handlers"
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
}

func main() {
	handler()
	// lambda.Start(handler)
}
