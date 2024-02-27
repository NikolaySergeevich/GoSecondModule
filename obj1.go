package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func obj1() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nВведите число")
		stdIn, _ := reader.ReadString('\n')
		stdIn = strings.TrimSpace(stdIn)
		if len(stdIn) == 0 {
			fmt.Println("Нельзя вводить пустое значение")
			continue
		}
		num, er := strconv.Atoi(stdIn)
		if er != nil {
			stdIn = stdIn[:len(stdIn)-1]
			if stdIn == "exit" {
				break
			} else {
				fmt.Println("Вы ввели не число")
				continue
			}
		}
		fmt.Println("Ваше число - ", num)
		chf := square(num)
		fmt.Println("произвдение = ", <-mult(chf))
	}
}

func square(num int) <-chan int {
	ch := make(chan int)
	go func() {
		fmt.Println("Квадрат числа = ", num * num)
		ch <- num * num
	}()
	return ch
}

func mult(ch <-chan int) <-chan int {
	newCh := make(chan int)
	go func() {
		num := <-ch
		newCh <- num * 2
	}()
	return newCh
}
