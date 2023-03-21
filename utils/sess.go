package utils

import "github.com/aws/aws-sdk-go/aws/session"

func Sess() *session.Session {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return sess
}
