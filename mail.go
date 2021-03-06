package main

import (
	"log"
	"net/smtp"
	"strings"
)

// SendMail is the function to send the mails
func SendMail(user, password, host, to, subject, body, mailtype string) error {
	log.Println("Sent to " + to)
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var contentType string
	if mailtype == "html" {
		contentType = "Content-Type: text/html;charset=UTF-8"
	} else {
		contentType = "Content-Type: text/plain;charset=UTF-8"
	}
	body = strings.TrimSpace(body)
	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	sentTo := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, sentTo, msg)
	if err != nil {
		log.Println(err.Error())
	}
	return err
}
