package ec2

import (
	"github.com/MZCBBD/AWSScheduler/utils"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2 struct {
	Id     string
	Status string
	Msg    string
}

// instanceIDs := []*string {
// 	// mzc-bbd2 vault-1
// 	aws.String("i-0420a345119580384"),
// 	// mzc-bbd2 vault-2
// 	aws.String("i-03ecc4d66213d0b37"),
// }

func New(id, status, msg string) *EC2 {
	return &EC2{Id: id, Status: status, Msg: msg}
}

func (e *EC2) GetStatus() (string, error) {
	sess := utils.Sess()
	svc := ec2.New(sess)

	return "sample", nil
}
