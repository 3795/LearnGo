package main

import (
	"Project/LearnGo/Learn/algorithm/array/array_util"
	"fmt"
)

/**
使用递归法实现二分搜索
*/
func main() {
	arr := array_util.GenArr(10)
	arr[4] = 68
	fmt.Printf("数组为 %v \n", arr)
	index := binarySearch2(arr, 0, len(arr)-1, 68)
	fmt.Printf("查找的索引为 %d", index)
}

func binarySearch2(arr []int, left, right, target int) int {
	if left > right {
		return -1
	}

	middle := left + (right-left)/2
	if target == arr[middle] {
		return middle
	} else if target > arr[middle] {
		return binarySearch2(arr, left, middle-1, target)
	} else {
		return binarySearch2(arr, middle+1, right, target)
	}
}
