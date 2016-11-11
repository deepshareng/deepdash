package email

import (
	"strings"

	"github.com/qiniu/log"

	"gopkg.in/gomail.v2"
)

func NewEmail(to, subject, msg, format string) *Email {
	return &Email{to: to, subject: subject, msg: msg, format: format}
}

func DeepshareDefaultEmailHost() *HostMail {
	return &HostMail{User: USER, Password: PASSWORD, Host: HOST}
}

func (email *Email) SendEmail(hostMail *HostMail) error {
	log.Info("[Email] Try to mail " + email.subject)

	m := gomail.NewMessage()
	toEmails := strings.Split(email.to, ",")
	m.SetHeader("From", hostMail.User)
	m.SetHeader("To", toEmails...)
	m.SetHeader("Subject", email.subject)
	m.SetBody(email.format, email.msg)
	d := gomail.NewPlainDialer(hostMail.Host, 25, hostMail.User, hostMail.Password)

	if err := d.DialAndSend(m); err != nil {
		log.Info("[Email] error ", err)
		return err
	}
	log.Info("[Email] send " + email.subject + " success!")
	return nil
}
