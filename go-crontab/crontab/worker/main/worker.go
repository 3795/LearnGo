package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan int)
	timer := time.NewTimer(1 * time.Second)
	go func() {
		for {
			ch1 <- 1
			time.Sleep(1 * time.Second)
		}
	}()
	for {
		select {
		case a := <-ch1:
			fmt.Println(a)
		case <-timer.C:
			fmt.Println("超时了")
		}
		timer.Reset(1 * time.Second)
	}
	fmt.Println("你好")
}
