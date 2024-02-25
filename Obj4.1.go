package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func Obj41(fileName string){
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			newF(fileName)
		}
	} else {
		oldFile(fileName)
	}
}

func newF(fileName string) error{
	var b bytes.Buffer
	count := 1
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
		b.WriteString(fmt.Sprintf("%d. %v-%v-%v %v:%v:%v - %s\n", count, date.Year(), date.Month(), date.Day(), date.Hour(),
									date.Minute(), date.Second(), message))
		count = count + 1
	}
	if err := ioutil.WriteFile(fileName, b.Bytes(), 0666); err != nil{
		return errors.New("не смог записать в файл")
	}
	return nil
}

