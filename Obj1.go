package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func Obj1(fileName string) {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			newFile(fileName)
		}
	} else {
		oldFile(fileName)
	}
}

func newFile(fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		return
	}
	defer file.Close()
	reader := bufio.NewReader(os.Stdin)
	countStr := 1
	for {
		fmt.Println("Введите ваше сообщение")
		message, _ := reader.ReadString('\n')
		num := strings.LastIndex(message, "\n")
		message = message[:num-1]
		if message == "exit" {
			break
		}
		date := time.Now()
		file.WriteString(fmt.Sprintf("%d. %v-%v-%v %v:%v:%v - %s\n", countStr, date.Year(), date.Month(), date.Day(), date.Hour(),
			date.Minute(), date.Second(), message))
		countStr = countStr + 1
	}
}

func oldFile(fileName string) {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	countStr := countSTR(fileName) + 1
	defer file.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Введите ваше сообщение")
		message, _ := reader.ReadString('\n')
		num := strings.LastIndex(message, "\n")
		message = message[:num-1]
		if message == "exit" {
			break
		}
		date := time.Now()
		file.WriteString(fmt.Sprintf("%d. %v-%v-%v %v:%v:%v - %s\n", countStr, date.Year(), date.Month(), date.Day(), date.Hour(),
			date.Minute(), date.Second(), message))
		countStr = countStr + 1
	}

}

func countSTR(fileName string) int {
	file, er := os.Open(fileName)
	if er != nil {
		return 0
	}
	fileScanner := bufio.NewScanner(file)
	countStr := 0
	for fileScanner.Scan() {
		countStr = countStr + 1
	}
	return countStr
}
