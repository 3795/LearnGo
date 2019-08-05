package main

import "fmt"

type MinStack struct {
	data []int
	helper []int
}

func Constructor() MinStack {
	ms := MinStack{}
	return ms
}

func (this *MinStack) Push(x int) {
	this.data = append(this.data, x)
	if len(this.helper) == 0 || x <= this.helper[len(this.helper)-1] {
		this.helper = append(this.helper, x)
	} else {
		this.helper = append(this.helper, this.helper[len(this.helper)-1])
	}
}

func (this *MinStack) Pop() {
	if len(this.data) == 0 || len(this.helper) == 0 {
		panic("非法操作")
	}
	this.data = this.data[:len(this.data)-1]
	this.helper = this.helper[:len(this.helper)-1]
}

func (this *MinStack) Top() int {
	if len(this.data) == 0 {
		panic("非法操作")
	}
	return this.data[len(this.data)-1]
}

func (this *MinStack) GetMin() int {
	if len(this.helper) == 0 {
		panic("操作非法")
	}
	return this.helper[len(this.helper)-1]
}

func main() {
	ms := Constructor()
	ms.Push(3)
	fmt.Println(ms.GetMin())
	fmt.Println(ms.Top())
}