package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func obj2() {

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	prNum := make(chan int)

	n := 0
	for {
		n = n + 1
		if n <= 5 {
			go func() {
				prNum <- n
			}()
		}

		time.Sleep(time.Second)
		select {
		case <-exit:
			fmt.Println("Корректно закрыл программу")
			return
		case a := <-prNum:
			fmt.Println(a)
		default:
			fmt.Println("ку-ку")

		}
	}
}
