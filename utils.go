package main

import (
	"bufio"
	"encoding/gob"
	"errors"
	"io"
	"log"
	"os"
	"regexp"
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

func delArrayVar(arr []string, str string) []string {
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

func loadData(ct *Content, path string) {
	log.Println("Load Data...")
	file, err := os.Open(path)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer file.Close()
	dec := gob.NewDecoder(file)
	err = dec.Decode(ct)
	if err != nil {
		log.Println(err.Error())
	}
}

func saveData(ct Content, path string) {
	log.Println("Save Data...")
	file, err := os.Create(path)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer file.Close()
	enc := gob.NewEncoder(file)
	err = enc.Encode(ct)
	if err != nil {
		log.Println(err.Error())
	}
}

func emailsValid(emails []string) []string {
	res := make([]string, 0)
	for _, email := range emails {
		if valid, err := regexp.Match("^[\\w.-]+@[\\w-]+(.[\\w-]+)*.[\\w]{2,6}$",
			[]byte(email)); valid && err == nil {
			res = append(res, email)
		} else {
			log.Println("Drop:", email)
		}
	}
	return res
}
