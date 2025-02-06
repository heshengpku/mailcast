package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	runing   bool
	mailType = true
	data     = "data.dat"
	ct       Content
)

type Content struct {
	Name     string   // 发件人邮箱
	Pwd      string   // 发件人密码
	Host     string   // SMTP服务器地址
	MailList []string // 收件人列表
	Subject  string   // 邮件主题
	Body     string   // 邮件正文
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("邮件发送工具")

	// 创建多行文本输入框和按钮
	contentEntry := widget.NewMultiLineEntry()
	sendButton := widget.NewButton("发送", func() {
		// TODO: 实现发送逻辑
	})
	clearButton := widget.NewButton("清除", func() {
		contentEntry.SetText("")
	})

	// 创建水平布局容器
	buttonBox := container.NewHBox(sendButton, clearButton)
	
	// 创建垂直布局容器
	mainContainer := container.NewVBox(
		contentEntry,
		buttonBox,
	)

	myWindow.SetContent(mainContainer)
	myWindow.Resize(fyne.NewSize(600, 400))
	myWindow.ShowAndRun()
}

func sendThread(msgbox, es *widget.Entry, progress *widget.ProgressBar) {
	mails, err := readMailsFromFile(data)
	if err != nil {
		msgbox.Append("读取邮件列表失败：" + err.Error() + "\n")
		return
	}
	
	count := len(mails)
	success := 0
	msgbox.Append(">>>开始发送，共" + strconv.Itoa(count) + "条\n\n")
	startT := time.Now()
	
	for index, to := range mails {
		if !runing {
			break
		}
		
		msgbox.Append("发送到 " + to + " ... " + strconv.Itoa((index+1)*100/count) + "%\n")
		progress.SetValue(float64(index + 1) / float64(count))
		
		var err error
		if mailType {
			err = SendMail(ct.Name, ct.Pwd, ct.Host, to, ct.Subject, ct.Body, "html")
		} else {
			err = SendMail(ct.Name, ct.Pwd, ct.Host, to, ct.Subject, ct.Body, "plain")
		}
		
		if err != nil {
			msgbox.Append("发送失败：" + err.Error() + "\n\n")
			mails = delMailFromList(mails, to)
			es.SetText(strings.Join(mails, "\r\n"))
			time.Sleep(time.Second)
			continue
		}
		
		success++
		msgbox.Append("发送成功！\n\n")
		mails = delMailFromList(mails, to)
		es.SetText(strings.Join(mails, "\r\n"))
		time.Sleep(100 * time.Millisecond)
	}
	
	duration := time.Since(startT)
	msgbox.Append(fmt.Sprintf("<<<停止发送！成功%d条，用时%v\n\n", success, duration))
}
