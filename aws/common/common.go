package common

type Handler interface {
	GetStatus()
	Stop()
	Start()
}
