package mailer

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"ta-spbe-backend/config"
	"ta-spbe-backend/service"
	"text/template"
)

type SimpleMailer struct {
	Auth  smtp.Auth
	Debug bool
	Host  string
	Port  int
	From  string
	body  string
}

const (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

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

func (m *SimpleMailer) Send(subject, message []byte, receiver []string, templateName string, items interface{}) error {
	// toSend := "Subject: " + string(subject) + "\n\n" + string(message)
	err := m.parseTemplate(templateName, items)
	if err != nil {
		log.Println(err)
	}

	toSend := "Subject: " + string(subject) + "\r\n" + MIME + "\r\n" + m.body
	if m.Debug {
		log.Println(toSend)
		return nil
	}

	return smtp.SendMail(fmt.Sprintf("%s:%d", m.Host, m.Port), m.Auth, m.From, receiver, []byte(toSend))
}

func (r *SimpleMailer) parseTemplate(fileName string, data interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}
	r.body = buffer.String()
	return nil
}
