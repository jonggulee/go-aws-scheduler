package common

type Handler interface {
	GetStatus() error
	Stop() (string, error)
	Start() (string, error)
}
