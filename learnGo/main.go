package main

import (
	"errors"
	"fmt"
)

var v = "Test"

func main() {
	outerFunc()
}

func outerFunc() {
	defer func() {
		p := recover()
		if p !=  nil {
			fmt.Println(p)
		}
	}()
	innerFunc()
}

func innerFunc() {
	panic(errors.New("出现了一个Panic"))
}
