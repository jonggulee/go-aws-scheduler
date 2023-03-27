package aws

import (
	"github.com/MZCBBD/AWSScheduler/aws/asg"
	"github.com/MZCBBD/AWSScheduler/aws/common"
	"github.com/MZCBBD/AWSScheduler/aws/ec2"
	"github.com/MZCBBD/AWSScheduler/aws/rds"
)

func NewAwsScheduler(service string, Id string) common.Handler {
	switch service {
	case "ec2":
		return ec2.New(Id, "", "", false)
	case "rds":
		return rds.New(Id, "", "", false)
	case "asg":
		return asg.New(Id, "", 0, false)
	default:
		return ec2.New(Id, "", "", false)
	}
}
