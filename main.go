package main

import (
	"os"

	"github.com/MZCBBD/AWSScheduler/rds"
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
