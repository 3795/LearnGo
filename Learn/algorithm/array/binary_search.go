package main

import (
	"Project/LearnGo/Learn/algorithm/array/array_util"
	"fmt"
)

/**
二分查找法实现
*/
func main() {
	arr := array_util.GenArr(10)
	arr[4] = 68
	fmt.Printf("数组为 %v \n", arr)
	index := binarySearch(arr, 88)
	fmt.Printf("查找的索引为 %d", index)
}

func binarySearch(arr []int, target int) int {
	l := 0
	r := len(arr) - 1
	var middle int
	for {
		if l > r {
			break
		}
		middle = l + (r-l)/2
		if target == arr[middle] {
			return middle
		} else if target < arr[middle] {
			r = middle - 1
		} else {
			l = middle + 1
		}
	}
	return -1
}
