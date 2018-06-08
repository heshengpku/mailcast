package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
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
	var emails, body, msgbox *walk.TextEdit
	var user, password, host, subject *walk.LineEdit
	var readBtn, startBtn *walk.PushButton
	mw := &MyMainWindow{}
	MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "邮件群发器",
		MinSize:  Size{800, 600},
		Layout:   HBox{},
		Children: []Widget{
			VSplitter{
				Children: []Widget{
					TextEdit{AssignTo: &emails, Text: ct.Send, ToolTipText: "待发送邮件列表，每列一个"},
					PushButton{
						AssignTo:  &readBtn,
						Text:      "打开",
						OnClicked: mw.pbClicked,
					},
				},
			},
			VSplitter{
				Children: []Widget{
					LineEdit{AssignTo: &user, Text: ct.Name, CueBanner: "请输入邮箱用户名"},
					LineEdit{AssignTo: &password, Text: ct.Pwd, PasswordMode: true, CueBanner: "请输入邮箱登录密码"},
					LineEdit{AssignTo: &host, Text: ct.Host, CueBanner: "SMTP服务器:端口"},
					LineEdit{AssignTo: &subject, Text: ct.Subject, CueBanner: "请输入邮件主题……"},
					TextEdit{AssignTo: &body, Text: ct.Body, ToolTipText: "请输入邮件内容", ColumnSpan: 2},
					TextEdit{AssignTo: &msgbox, ReadOnly: true},
					PushButton{
						AssignTo: &startBtn,
						Text:     "开始群发",
						OnClicked: func() {
							ct.Name = user.Text()
							ct.Pwd = password.Text()
							ct.Host = host.Text()
							ct.Subject = subject.Text()
							ct.Body = body.Text()
							ct.Send = emails.Text()
							SaveData()

							if runing == false {
								runing = true
								startBtn.SetText("停止发送")
								go sendThread(msgbox, emails)
							} else {
								runing = false
								startBtn.SetText("开始群发")
							}
						},
					},
				},
			},
		},
	}.Run()
}

func sendThread(msgbox, es *walk.TextEdit) {
	sentTo := strings.Split(ct.Send, "\r\n")
	count := len(sentTo)
	success := 0
	for index, to := range sentTo {
		if runing == false {
			break
		}
		msgbox.SetText("发送到" + to + "..." + strconv.Itoa((index/count)*100) + "%")
		err := SendMail(ct.Name, ct.Pwd, ct.Host, to, ct.Subject, ct.Body, "html")
		if err != nil {
			msgbox.AppendText("\r\n失败：" + err.Error() + "\r\n")
			if err.Error() == "550 Mailbox not found or access denied" {
				ct.Send = strings.Join(DelArrayVar(strings.Split(ct.Send, "\r\n"), to), "\r\n")
				es.SetText(ct.Send)
			}
			time.Sleep(1 * time.Second)
			continue
		} else {
			success++
			msgbox.AppendText("\r\n发送成功！")
			ct.Send = strings.Join(DelArrayVar(strings.Split(ct.Send, "\r\n"), to), "\r\n")
			es.SetText(ct.Send)
		}
		time.Sleep(1 * time.Second)
	}
	SaveData()
	msgbox.AppendText("停止发送！成功" + strconv.Itoa(success) + "条\r\n")
}

type MyMainWindow struct {
	*walk.MainWindow
	edit *walk.TextEdit

	path string
}

func (mw *MyMainWindow) pbClicked() {
	dlg := new(walk.FileDialog)
	dlg.FilePath = mw.path
	dlg.Title = "Select File"
	dlg.Filter = "Txt files (*.txt)|*.txt|All files (*.*)|*.*"

	if ok, err := dlg.ShowOpen(mw); err != nil {
		mw.edit.AppendText("Error: file open fail\r\n")
		return
	} else if !ok {
		mw.edit.AppendText("Cancel\r\n")
		return
	}
	mw.path = dlg.FilePath
	mw.edit.AppendText(fmt.Sprintf("Select : %s\r\n", mw.path))
}
