package main

import (
	"fmt"
	"time"
)

func main() {

	fmt.Println("Hello World")

	time.Sleep(time.Duration(1) * time.Minute)

	fmt.Println("睡眠结束")
}
