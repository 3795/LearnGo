package main

import (
	"fmt"
	"sync"
)

/**
交替打印FooBar
https://leetcode-cn.com/problems/print-foobar-alternately/
*/
func main() {
	executor(100)
}

func executor(n int) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	fooChan := make(chan int) // 防止最后一个元素阻塞通道
	barChan := make(chan int)
	go foo(&wg, n, fooChan, barChan)
	go bar(&wg, n, fooChan, barChan)
	fooChan <- 1
	wg.Wait()
}

func foo(wg *sync.WaitGroup, n int, fooCh chan int, barCh chan int) {
	for i := 0; i < n; i++ {
		select {
		case <-fooCh:
			fmt.Println("foo")
			barCh <- 1
		}
	}
	wg.Done()
}

func bar(wg *sync.WaitGroup, n int, fooCh chan int, barCh chan int) {
	for i := 0; i < n; i++ {
		select {
		case <-barCh:
			fmt.Println("bar")
			// 防止最后一个元素阻塞通道，导致协程被阻塞
			if i < 99 {
				fooCh <- 1
			}
		}
	}
	wg.Done()
}
