package main

import (
	"os"

	"github.com/MZCBBD/AWSScheduler/ec2"
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
	if os.Getenv("target") == "ec2" {
		// if os.Getenv("env") == "stop" {
		// ec2.StopDBInstanceHandler()
		// }
		if os.Getenv("env") == "start" {
			ec2.StartInstanceHandler()
		}
	}
}

func main() {
	handler()
	// lambda.Start(handler)
}
