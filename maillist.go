package main

import "errors"

// Content holds the content struct
type Content struct {
	Name    string
	Pwd     string
	Host    string
	Subject string
	Body    string
}

// ct is the email configure
var ct Content

func load(path string) error {
	return loadData(&ct, path)
}

func save(path string) error {
	return saveData(ct, path)
}

func setContent(name, pwd, host, subject, body string) error {
	if emailValid(name) == false {
		return errors.New("Invalid email format")
	}
	ct.Name = name
	ct.Pwd = pwd
	ct.Host = host
	ct.Subject = subject
	ct.Body = body
	return nil
}

// mailList is the mail list
var mailList []string

func loadMails(path string) error {
	mails, err := readMailsFromFile(path)
	if err != nil {
		return err
	}
	mailList = make([]string, len(mails))
	mailList = mails
	return nil
}

func getMails() []string {
	return mailList
}

func delMail(mail string) {
	mailList = delMailFromList(mailList, mail)
}
