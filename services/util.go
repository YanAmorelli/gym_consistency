package services

import (
	"crypto/tls"
	"math/rand"
	"time"

	"github.com/labstack/gommon/log"
	gomail "gopkg.in/mail.v2"
)

func GenerateRandomPassword() string {
	rand.Seed(int64(time.Now().Unix()))
	length := 24
	ran_str := make([]byte, length)

	for i := 0; i < length; i++ {
		ran_str[i] = byte(65 + rand.Intn(26))
	}

	password := string(ran_str)

	for i := 0; i < 3; i++ {
		str := []byte(password)
		position := rand.Intn(23)
		str[position] = '@'
		password = string(str)
	}

	return password
}

func SendEmail(fromEmail, password, toEmail, subjectEmail, bodyEmail string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fromEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subjectEmail)
	m.SetBody("text/plain", bodyEmail)
	d := gomail.NewDialer("smtp.gmail.com", 587, fromEmail, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		log.Error("error trying to send email. Error: ", err.Error())
		return err
	}

	return nil
}
