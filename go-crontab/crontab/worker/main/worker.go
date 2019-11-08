package main

import "fmt"

func main() {
	ch1 := make(chan int, 10)
	for i := 0; i < 11; i++ {
		append(ch1, i)
	}
	fmt.Println("程序未被阻塞")
}

func append(ch chan int, i int) {
	select {
	case ch <- i:
	default:
		fmt.Printf("数字 %d 被丢弃了\n", i)
	}
}
