package mailer

import (
	"fmt"
	"log"
	"net/smtp"
	"ta-spbe-backend/config"
	"ta-spbe-backend/service"
)

type SimpleMailer struct {
	Auth  smtp.Auth
	Debug bool
	Host  string
	Port  int
	From  string
}

func NewSimpleMailer(smtpCfg config.SMTPClient) (service.Mailer, error) {
	return &SimpleMailer{
		Auth: smtp.PlainAuth(
			smtpCfg.AdminIdentity,
			smtpCfg.AdminEmail,
			smtpCfg.AdminPassword,
			smtpCfg.Host,
		),
		Debug: smtpCfg.Debug,
		Host:  smtpCfg.Host,
		Port:  smtpCfg.Port,
		From:  smtpCfg.AdminEmail,
	}, nil
}

func (m *SimpleMailer) Send(subject, message []byte, receiver []string) error {
	toSend := "Subject: " + string(subject) + "\n\n" + string(message)
	if m.Debug {
		log.Println(toSend)
		return nil
	}
	return smtp.SendMail(fmt.Sprintf("%s:%d", m.Host, m.Port), m.Auth, m.From, receiver, []byte(toSend))
}
