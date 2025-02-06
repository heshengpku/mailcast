package main

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)


// readMailsFromFile 从指定文件中读取邮件列表
// filename: 待读取的文件路径
// 返回值: 邮件列表切片和可能的错误
func readMailsFromFile(filename string) ([]string, error) {
	result := make([]string, 0)
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}
	
	return result, nil
}


// readTxtFromFile 读取文本文件的全部内容
// filename: 待读取的文件路径
// 返回值: 文件内容字符串和可能的错误
func readTxtFromFile(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}
	return string(content), nil
}

// delMailFromList 从邮件列表中删除指定邮件地址
// arr: 原邮件列表
// str: 要删除的邮件地址
// 返回值: 删除后的邮件列表
func delMailFromList(arr []string, str string) []string {
	str = strings.TrimSpace(str)
	for i, v := range arr {
		if v == str {
			return append(arr[:i], arr[i+1:]...)
		}
	}
	return arr
}

// emailValid 验证邮件地址格式是否合法
// email: 待验证的邮件地址
// 返回值: 是否合法
func emailValid(email string) bool {
	pattern := `^[\w.-]+@[\w-]+(\.[\w-]+)*\.[a-zA-Z]{2,6}$`
	valid, err := regexp.MatchString(pattern, email)
	return err == nil && valid
}

func loadData(ct *Content, path string) error {
	log.Println("Load Data...")
	file, err := os.Open(path)
	if err != nil {
		log.Println(err.Error())
		return errors.New("Failed to open " + path + ":" + err.Error())
	}
	defer file.Close()
	dec := gob.NewDecoder(file)
	err = dec.Decode(ct)
	if err != nil {
		log.Println(err.Error())
		return errors.New("Failed to decode from " + path + ":" + err.Error())
	}
	return nil
}

func saveData(ct Content, path string) error {
	log.Println("Save Data...")
	file, err := os.Create(path)
	if err != nil {
		log.Println(err.Error())
		return errors.New("Failed to open " + path + ":" + err.Error())
	}
	defer file.Close()
	enc := gob.NewEncoder(file)
	err = enc.Encode(ct)
	if err != nil {
		log.Println(err.Error())
		return errors.New("Failed to encode to " + path + ":" + err.Error())
	}
	return nil
}
