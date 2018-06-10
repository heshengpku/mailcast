package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ProtonMail/ui"
)

type Content struct {
	Name    string
	Pwd     string
	Host    string
	Subject string
	Body    string
	Send    string
}

var ct Content
var runing bool

func main() {
	LoadData()

	err := ui.Main(func() {
		emails := ui.NewMultilineEntry()
		openfile := ui.NewButton("选择文件")
		send := ui.NewButton("开始群发")

		vb1 := ui.NewVerticalBox()
		vb1.SetPadded(true)
		vb1hbox := ui.NewHorizontalBox()
		vb1hbox.Append(ui.NewLabel("邮箱列表"), false)
		vb1hbox.Append(openfile, false)
		vb1hbox.Append(ui.NewLabel(""), true)
		vb1.Append(vb1hbox, false)
		vb1.Append(emails, true)
		vb1.Append(send, false)

		user := ui.NewEntry()
		password := ui.NewEntry()
		host := ui.NewEntry()
		port := ui.NewEntry()
		subject := ui.NewEntry()
		body := ui.NewMultilineEntry()
		openbody := ui.NewButton("选择文件")

		vb2 := ui.NewVerticalBox()
		vb2.Append(ui.NewLabel("请输入发送邮箱用户名"), false)
		vb2.Append(user, false)
		vb2.Append(ui.NewLabel("请输入发送邮箱登录密码"), false)
		vb2.Append(password, false)
		vb2hbox := ui.NewHorizontalBox()
		vb2hbox.Append(ui.NewLabel("SMTP服务器"), false)
		vb2hbox.Append(host, true)
		vb2hbox.Append(ui.NewLabel("端口"), false)
		vb2hbox.Append(port, true)
		vb2.Append(vb2hbox, false)
		vb2.Append(ui.NewLabel("请输入邮件主题……"), false)
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

		window := ui.NewWindow("邮件群发器", 1000, 800, false)
		window.SetMargined(true)
		window.SetChild(hbox)

		openfile.OnClicked(func(*ui.Button) {
			path := ui.OpenFile(window)
			emailsText, err := readLine2Array(path)
			if err != nil {
				msgbox.Append(fmt.Sprintf("Select : %s\r\nerror: %s", path, err.Error()))
			}
			emails.SetText(strings.Join(emailsText, "\r\n"))
		})

		openbody.OnClicked(func(*ui.Button) {
			path := ui.OpenFile(window)
			bodyText, err := readLine2Array(path)
			if err != nil {
				msgbox.Append(fmt.Sprintf("Select : %s\r\nerror: %s", path, err.Error()))
			}
			body.SetText(strings.Join(bodyText, "\r\n"))
		})

		send.OnClicked(func(*ui.Button) {
			ct.Name = user.Text()
			ct.Pwd = password.Text()
			ct.Host = host.Text() + ":" + port.Text()
			ct.Subject = subject.Text()
			ct.Body = body.Text()
			ct.Send = emails.Text()
			SaveData()

			if runing == false {
				runing = true
				send.SetText("停止发送")
				go sendThread(msgbox, emails, progress)
			} else {
				runing = false
				send.SetText("开始群发")
			}

		})

		window.OnClosing(func(*ui.Window) bool {
			SaveData()
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
	sentTo := strings.Split(ct.Send, "\r\n")
	count := len(sentTo)
	success := 0
	for index, to := range sentTo {
		if runing == false {
			break
		}
		msgbox.Append("发送到" + to + "..." + strconv.Itoa((index/count)*100) + "%")
		progress.SetValue((index / count) * 100)
		err := SendMail(ct.Name, ct.Pwd, ct.Host, to, ct.Subject, ct.Body, "html")
		if err != nil {
			msgbox.Append("\r\n失败：" + err.Error() + "\r\n")
			if err.Error() == "550 Mailbox not found or access denied" {
				ct.Send = strings.Join(DelArrayVar(strings.Split(ct.Send, "\r\n"), to), "\r\n")
				es.SetText(ct.Send)
			}
			time.Sleep(1 * time.Second)
			continue
		} else {
			success++
			msgbox.Append("\r\n发送成功！")
			ct.Send = strings.Join(DelArrayVar(strings.Split(ct.Send, "\r\n"), to), "\r\n")
			es.SetText(ct.Send)
		}
		time.Sleep(1 * time.Second)
	}
	SaveData()
	msgbox.Append("停止发送！成功" + strconv.Itoa(success) + "条\r\n")
}
