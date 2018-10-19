package main

import (
	"bufio"
	"encoding/gob"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

// type Mail struct {
// 	Email string `json:"email"`
// 	Name  string `json:"name"`
// }

// func (mail Mail) String() string {
// 	return fmt.Sprintf("%s \t %s", mail.Email, mail.Name)
// }

// func mailArrayToStrings(list []Mail) []string {
// 	var strArray []string
// 	for _, email := range mailList {
// 		strArray = append(strArray, email.String())
// 	}
// 	return strArray
// }

func readMailsFromFile(filename string) ([]string, error) {
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

// func readMailsFromFile(filename string) ([]Mail, error) {
// 	result := make([]Mail, 0)
// 	file, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		return result, errors.New("Open file failed")
// 	}
// 	r := csv.NewReader(strings.NewReader(string(file[:])))
// 	records, err := r.Read()
// 	if err != nil {
// 		if err != io.EOF {
// 			return result, errors.New("ReadLine not finish")
// 		}
// 	}
// 	for _, line := range records {
// 		record := strings.Split(line, ",")
// 		if emailValid(record[0]) {
// 			mail := Mail{Email: record[0], Name: record[1]}
// 			result = append(result, mail)
// 		}
// 	}
// 	return result, nil
// }

func readTxtFromFile(filename string) (string, error) {
	context, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(context), nil
}

func delMailFromList(arr []string, str string) []string {
	str = strings.TrimSpace(str)
	for i, v := range arr {
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

func emailValid(email string) bool {
	if valid, err := regexp.Match("^[\\w.-]+@[\\w-]+(.[\\w-]+)*.[\\w]{2,6}$",
		[]byte(email)); valid && err == nil {
		return true
	}

	return false
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
