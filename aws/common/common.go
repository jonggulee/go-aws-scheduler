package common

import "github.com/aws/aws-sdk-go/aws/session"

type Handler interface {
	GetStatus()
	Stop()
	Start()
}

func Sess() *session.Session {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return sess
}
