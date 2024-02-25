package main

import (
	"io/ioutil"
	"os"
)

func Obj42(fileName string) (res string){
	file, err := os.Open(fileName)
	if err != nil {
		return "Не смог прочитать такой файл"
	}
	resB, er := ioutil.ReadAll(file)
	if er != nil {
		return "Не смог считать информацию из файла"
	}
	return string(resB)
}