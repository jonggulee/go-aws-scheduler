package common

type Handler interface {
	GetStatus() (string, error)
	// Stop(string) (string, error)
	// Start()
}
