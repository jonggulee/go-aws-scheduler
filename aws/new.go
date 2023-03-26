package aws

import (
	"github.com/MZCBBD/AWSScheduler/aws/common"
	"github.com/MZCBBD/AWSScheduler/aws/ec2"
)

func NewAwsScheduler(service string, Id string) common.Handler {
	switch service {
	case "ec2":
		return ec2.New(Id, "", "")
	case "rds":
		// 	return rds.New("", "", "")
		return nil
	default:
		return ec2.New(Id, "", "")
	}
}
