package main

import (
	"fmt"
)

func main(){
	fmt.Println("Задание 1")
	Obj1("Task1.txt")
	fmt.Println("=========================")
	fmt.Println("задание 2")
	fmt.Println(Obj2("Task1.txt"))
	fmt.Println("=========================")
	fmt.Println("задание 3")
	fmt.Println(Obj3("Task3.txt"))
	fmt.Println("=========================")
	fmt.Println("задание 4.1")
	Obj41("Task4.txt")
	fmt.Println("=========================")
	fmt.Println("задание 4.2")
	fmt.Println(Obj42("Task4.txt"))	
}