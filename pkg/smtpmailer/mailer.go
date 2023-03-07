package smtpmailer

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"text/template"
)

type Config struct {
	Debug         bool   `toml:"debug"`
	Host          string `toml:"host"`
	Port          int    `toml:"port"`
	AdminIdentity string `toml:"admin_identity"`
	AdminEmail    string `toml:"admin_email"`
	AdminPassword string `toml:"admin_password"`
}

type Client interface {
	Send(subject, message []byte, receiver []string, templateName string, items interface{}) error
}

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

func NewSimpleMailer(smtpCfg Config) (Client, error) {
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
