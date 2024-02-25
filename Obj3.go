package main

import (
	"errors"
	"fmt"
	"io"
	"os"

)

func Obj3(fileName string) string {
	if er := createF(fileName); er != nil {
		return fmt.Sprintf("%v", er)
	}
	str, err := readF(fileName)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	return str	
}

func createF(fileName string) error {
	byteStr := []byte("Этот текст для проверки.")
	if err := os.WriteFile(fileName, byteStr, 0444); err != nil {
		return err
	}
	return nil
}

func readF(fileName string) (_ string, e error) {
	file, er := os.Open(fileName)
	if er != nil {
		e = errors.New("не смог открыть файл")
		return 
	}
	defer file.Close()
	s, _ := os.Stat(fileName)
	fmt.Printf("Права доступа: %v\n",s.Mode())
	bufer := make([]byte, int(s.Size()))
	if _, err := io.ReadFull(file, bufer); err != nil{
		e = errors.New("не смог прочитать файл")
		return 
	}
	return string(bufer), nil
}

