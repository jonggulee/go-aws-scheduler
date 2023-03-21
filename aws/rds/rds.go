package rds

type Rds struct {
}

func New(id, status, msg string) *Rds {
	return &Rds{}
}

func (e *Rds) GetStatus() (string, error) {
	return "RDS sample", nil
}
