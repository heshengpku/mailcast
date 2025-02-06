package main

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

// 添加常量定义
const (
	htmlContentType = "Content-Type: text/html;charset=UTF-8"
	plainContentType = "Content-Type: text/plain;charset=UTF-8"
)

// SendMail is the function to send the mails
func SendMail(user, password, host, to, subject, body, mailtype string) error {
	// 参数验证
	if user == "" || password == "" || host == "" || to == "" {
		return fmt.Errorf("required parameters cannot be empty")
	}
	
	// 解析主机地址
	hp := strings.Split(host, ":")
	if len(hp) != 2 {
		return fmt.Errorf("invalid host format, expected host:port")
	}
	
	auth := smtp.PlainAuth("", user, password, hp[0])
	
	// 设置内容类型
	contentType := plainContentType
	if mailtype == "html" {
		contentType = htmlContentType
	}
	
	// 使用 strings.Builder 构建消息
	var msgBuilder strings.Builder
	msgBuilder.WriteString("To: " + to + "\r\n")
	msgBuilder.WriteString("From: " + user + "<" + user + ">\r\n")
	msgBuilder.WriteString("Subject: " + subject + "\r\n")
	msgBuilder.WriteString(contentType + "\r\n\r\n")
	msgBuilder.WriteString(strings.TrimSpace(body))
	
	msg := []byte(msgBuilder.String())
	sentTo := strings.Split(to, ";")
	
	// 发送邮件
	if err := smtp.SendMail(host, auth, user, sentTo, msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	
	// 发送成功后记录日志
	log.Printf("Successfully sent email to %s", to)
	return nil
}
