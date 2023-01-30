package service

type Mailer interface {
	Send(subject, message []byte, receiver []string) error
}
