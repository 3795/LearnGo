package main

import (
	"LearnGo/Learn/algorithm/array/array_util"
	"fmt"
)

/**
数组中的第K大的元素
*/
func main() {
	arr := array_util.GenArr(10)
	fmt.Printf("%v\n", arr)
	fmt.Println(findKthLargest(arr, 3))
}

func findKthLargest(nums []int, k int) int {
	sort(nums)
	return nums[len(nums)-k]
}

func partition(arr []int, left, right int) int {
	e := arr[left]
	j := left
	for i := left + 1; i < right; i++ {
		if arr[i] < e {
			j++
			swap2(arr, j, i)
		}
	}
	swap2(arr, left, j)
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

func swap2(arr []int, i, j int) {
	tmp := arr[i]
	arr[i] = arr[j]
	arr[j] = tmp
}
