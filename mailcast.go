package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ProtonMail/ui"
)

var runing bool
var mailType = true
var data = "data.dat"

func main() {
	load(data)

	err := ui.Main(func() {
		// 左栏
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

		// 中栏
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
		saveCt := ui.NewButton("保存设置")
		subject := ui.NewEntry()
		subject.SetText(ct.Subject)
		body := ui.NewMultilineEntry()
		openbody := ui.NewButton("选择文件")
		htmlCb := ui.NewCheckbox("html")
		htmlCb.SetChecked(true)

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
		vb2.Append(saveCt, false)
		vb2.Append(ui.NewLabel("请输入邮件主题:"), false)
		vb2.Append(subject, false)
		vb2hbox2 := ui.NewHorizontalBox()
		vb2hbox2.Append(ui.NewLabel("请输入邮件内容"), false)
		vb2hbox2.Append(openbody, false)
		vb2hbox2.Append(ui.NewLabel(""), true)
		vb2hbox2.Append(htmlCb, false)
		vb2.Append(vb2hbox2, false)
		vb2.Append(body, true)

		// 右栏
		msgbox := ui.NewMultilineEntry()
		msgbox.SetReadOnly(true)
		progress := ui.NewProgressBar()
		vb3 := ui.NewVerticalBox()
		vb3.Append(ui.NewLabel("日志"), false)
		vb3.Append(msgbox, true)
		vb3.Append(progress, false)

		// 三栏盒
		hbox := ui.NewHorizontalBox()
		hbox.Append(vb1, true)
		hbox.Append(vb2, true)
		hbox.Append(vb3, true)

		// 窗口
		window := ui.NewWindow("邮件群发器", 1280, 720, false)
		window.SetMargined(true)
		window.SetChild(hbox)

		// 按键控制
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

		htmlCb.OnToggled(func(*ui.Checkbox) {
			mailType = !mailType
		})

		openbody.OnClicked(func(*ui.Button) {
			path := ui.OpenFile(window)
			bodyText, err := readTxtFromFile(path)
			if err != nil {
				msgbox.Append(fmt.Sprintf("打开文件 %s\n错误: %s\n\n", path, err.Error()))
			}
			body.SetText(bodyText)
		})

		saveCt.OnClicked(func(b *ui.Button) {
			setContent(user.Text(),
				password.Text(),
				host.Text()+":"+port.Text(),
				subject.Text(),
				body.Text())
			err := save(data)
			if err != nil {
				msgbox.Append("保存失败！\n" + err.Error())
			} else {
				msgbox.Append("保存成功！\n")
			}
		})

		send.OnClicked(func(b *ui.Button) {
			setContent(user.Text(),
				password.Text(),
				host.Text()+":"+port.Text(),
				subject.Text(),
				body.Text())
			save(data)

			if runing == false {
				runing = true
				b.SetText("停止发送")
				vb2.Disable()
				go sendThread(msgbox, emails, progress)
			} else {
				runing = false
				b.SetText("开始群发")
				vb2.Enable()
			}
		})

		// 窗口控制
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
		var err error
		if mailType {
			msgbox.Append("True")
			// err = SendMail(ct.Name, ct.Pwd, ct.Host, to, ct.Subject, ct.Body, "html")
		} else {
			msgbox.Append("False")
			// err = SendMail(ct.Name, ct.Pwd, ct.Host, to, ct.Subject, ct.Body, "plain")
		}
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
		time.Sleep(100 * time.Millisecond)
	}
	duration := time.Since(startT)
	msgbox.Append(fmt.Sprintf("<<<停止发送！成功%d条，用时%v\n\n", success, duration))
}
