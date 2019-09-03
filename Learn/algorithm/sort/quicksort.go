package main

import (
	"LearnGo/Learn/algorithm/array/array_util"
	"fmt"
)

/**
快速排序实现
*/

func main() {
	arr := array_util.GenArr(10)
	fmt.Printf("%v\n", arr)
	sort(arr)
	fmt.Printf("%v\n", arr)
}

func partition(arr []int, left, right int) int {
	e := arr[left]
	j := left
	for i := left + 1; i < right; i++ {
		if arr[i] < e {
			j++
			array_util.Swap(arr, j, i)
		}
	}
	array_util.Swap(arr, left, j)
	return j
}

func sortDetail(arr []int, left, right int) {
	if left >= right {
		return
	}

	p := partition(arr, left, right)
	sortDetail(arr, left, p)
	sortDetail(arr, p+1, right)
}

func sort(arr []int) {
	sortDetail(arr, 0, len(arr))
}
