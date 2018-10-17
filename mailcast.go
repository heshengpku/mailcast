package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ProtonMail/ui"
)

var runing bool
var data string = "data.dat"

func main() {
	load(data)

	err := ui.Main(func() {
		emails := ui.NewMultilineEntry()
		openfile := ui.NewButton("选择文件")
		send := ui.NewButton("开始群发")

		vb1 := ui.NewVerticalBox()
		vb1.SetPadded(true)
		vb1hbox := ui.NewHorizontalBox()
		vb1hbox.Append(ui.NewLabel("群发邮箱列表，每行一个"), false)
		vb1hbox.Append(openfile, false)
		vb1hbox.Append(ui.NewLabel(""), true)
		vb1.Append(vb1hbox, false)
		vb1.Append(emails, true)
		vb1.Append(send, false)

		user := ui.NewEntry()
		user.SetText(ct.Name)
		password := ui.NewPasswordEntry()
		password.SetText(ct.Pwd)
		host := ui.NewEntry()
		strs := strings.Split(ct.Host, ":")
		if len(strs) >= 1 {
			host.SetText(strs[0])
		}
		port := ui.NewEntry()
		if len(strs) >= 2 {
			port.SetText(strs[1])
		}
		subject := ui.NewEntry()
		subject.SetText(ct.Subject)
		body := ui.NewMultilineEntry()
		openbody := ui.NewButton("选择文件")

		vb2 := ui.NewVerticalBox()
		vb2.Append(ui.NewLabel("发送邮箱:"), false)
		vb2.Append(user, false)
		vb2.Append(ui.NewLabel("登录密码:"), false)
		vb2.Append(password, false)
		vb2hbox := ui.NewHorizontalBox()
		vb2hbox.Append(ui.NewLabel("SMTP服务器:"), false)
		vb2hbox.Append(host, true)
		vb2hbox.Append(ui.NewLabel("端口:"), false)
		vb2hbox.Append(port, true)
		vb2.Append(vb2hbox, false)
		vb2.Append(ui.NewLabel("请输入邮件主题:"), false)
		vb2.Append(subject, false)
		vb2hbox2 := ui.NewHorizontalBox()
		vb2hbox2.Append(ui.NewLabel("请输入邮件内容"), false)
		vb2hbox2.Append(openbody, false)
		vb2hbox2.Append(ui.NewLabel(""), true)
		vb2.Append(vb2hbox2, false)
		vb2.Append(body, true)

		msgbox := ui.NewMultilineEntry()
		msgbox.SetReadOnly(true)
		progress := ui.NewProgressBar()
		vb3 := ui.NewVerticalBox()
		vb3.Append(ui.NewLabel("日志"), false)
		vb3.Append(msgbox, true)
		vb3.Append(progress, false)

		hbox := ui.NewHorizontalBox()
		hbox.Append(vb1, true)
		hbox.Append(vb2, true)
		hbox.Append(vb3, true)

		window := ui.NewWindow("邮件群发器", 1280, 720, false)
		window.SetMargined(true)
		window.SetChild(hbox)

		openfile.OnClicked(func(*ui.Button) {
			path := ui.OpenFile(window)
			err := loadMails(path)
			if err != nil {
				msgbox.Append(fmt.Sprintf("打开文件 %s\n错误: %s\n\n", path, err.Error()))
			} else {
				msgbox.Append(fmt.Sprintf("打开文件 %s\n成功\n\n", path))
			}
			emails.SetText(strings.Join(getMails(), "\r\n"))
		})

		openbody.OnClicked(func(*ui.Button) {
			path := ui.OpenFile(window)
			bodyText, err := readTxtFromFile(path)
			if err != nil {
				msgbox.Append(fmt.Sprintf("打开文件 %s\n错误: %s\n\n", path, err.Error()))
			}
			body.SetText(bodyText)
		})

		send.OnClicked(func(b *ui.Button) {
			ct.Name = user.Text()
			ct.Pwd = password.Text()
			ct.Host = host.Text() + ":" + port.Text()
			ct.Subject = subject.Text()
			ct.Body = body.Text()
			save(data)

			if runing == false {
				runing = true
				b.SetText("停止发送")
				go sendThread(msgbox, emails, progress)
			} else {
				runing = false
				b.SetText("开始群发")
			}

		})

		window.OnClosing(func(*ui.Window) bool {
			save(data)
			ui.Quit()
			return true
		})
		window.Show()
	})
	if err != nil {
		panic(err)
	}
}

func sendThread(msgbox, es *ui.MultilineEntry, progress *ui.ProgressBar) {
	mails := getMails()
	count := len(mails)
	success := 0
	msgbox.Append(">>>开始发送，共" + strconv.Itoa(count) + "条\n\n")
	startT := time.Now()
	for index, to := range mails {
		if runing == false {
			break
		}
		msgbox.Append("发送到 " + to + " ... " + strconv.Itoa((index+1)*100/count) + "%\n")
		progress.SetValue((index + 1) * 100 / count)
		err := SendMail(ct.Name, ct.Pwd, ct.Host, to, ct.Subject, ct.Body, "html")
		if err != nil {
			msgbox.Append("发送失败：" + err.Error() + "\n\n")
			delMail(to)
			es.SetText(strings.Join((getMails()), "\r\n"))
			time.Sleep(1 * time.Second)
			continue
		} else {
			success++
			msgbox.Append("发送成功！\n\n")
			delMail(to)
			es.SetText(strings.Join((getMails()), "\r\n"))
		}
		// time.Sleep(100 * time.Millisecond)
	}
	duration := time.Since(startT)
	msgbox.Append(fmt.Sprintf("<<<停止发送！成功%d条，用时%v\n\n", success, duration))
}
