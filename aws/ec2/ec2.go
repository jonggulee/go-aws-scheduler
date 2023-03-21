package ec2

type EC2 struct {
	Id     string
	Status string
	Msg    string
}

func New(id, status, msg string) *EC2 {
	return &EC2{Id: id, Status: status, Msg: msg}
}

func (e *EC2) GetStatus() (string, error) {
	return "sample", nil
}
