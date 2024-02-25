package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func Obj2(fileName string) string {
	size, er := sizeFile(fileName)
	if er != nil {
		return fmt.Sprintf("%v", er)
	}
	file, er := os.Open(fileName)
	if er != nil {
		return "Не смог открыть файл"
	}
	defer file.Close()
	bufer := make([]byte, size)
	if _, err := io.ReadFull(file, bufer); err != nil{
		return "Не смог прочитать файл"
	}
	return fmt.Sprintf("Вот содержимое файла %s:\n%s\n", fileName, bufer)
}

func sizeFile(fileName string) (res int, er error) {
	infF, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, errors.New("такого файла нет")
		}
	}
	return int(infF.Size()), nil
}
