package main

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func readLine2Array(filename string) ([]string, error) {
	result := make([]string, 0)
	file, err := os.Open(filename)
	if err != nil {
		return result, errors.New("Open file failed")
	}
	defer file.Close()
	bf := bufio.NewReader(file)
	for {
		line, isPrefix, err := bf.ReadLine()
		if err != nil {
			if err != io.EOF {
				return result, errors.New("ReadLine not finish")
			}
			break
		}
		if isPrefix {
			return result, errors.New("Line is too long")
		}
		str := string(line)
		result = append(result, str)
	}
	return result, nil
}

func DelArrayVar(arr []string, str string) []string {
	str = strings.TrimSpace(str)
	for i, v := range arr {
		v = strings.TrimSpace(v)
		if v == str {
			if i == len(arr) {
				return arr[0 : i-1]
			}
			if i == 0 {
				return arr[1:len(arr)]
			}
			return append(arr[0:i], arr[i+1:len(arr)]...)
		}
	}
	return arr
}

func LoadData() {
	fmt.Println("LoadData")
	file, err := os.Open("data.dat")
	defer file.Close()
	if err != nil {
		fmt.Println(err.Error())
		ct = Content{
			Name:    "用户名",
			Pwd:     "用户密码",
			Host:    "SMTP服务器:端口",
			Subject: "邮件主题",
			Body:    "邮件内容",
			Send:    "要发送的邮箱，每行一个",
		}
		return
	}
	dec := gob.NewDecoder(file)
	err = dec.Decode(&ct)
	if err != nil {
		fmt.Println(err.Error())
		ct = Content{
			Name:    "用户名",
			Pwd:     "用户密码",
			Host:    "SMTP服务器:端口",
			Subject: "邮件主题",
			Body:    "邮件内容",
			Send:    "要发送的邮箱，每行一个",
		}
	}
}

func SaveData() {
	fmt.Println("SaveData")
	file, err := os.Create("data.dat")
	defer file.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	enc := gob.NewEncoder(file)
	err = enc.Encode(ct)
	if err != nil {
		fmt.Println(err.Error())
	}
}
