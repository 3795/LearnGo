package main

import "fmt"

/**
颜色分类
https://leetcode-cn.com/problems/sort-colors/
使用计数排序的方法
*/

func main() {
	nums := []int{2, 0, 2, 1, 1, 0}
	sortColors(nums)
	fmt.Printf("%v", nums)
}

func sortColors(nums []int) {
	count := make([]int, 3)
	for i := 0; i < len(nums); i++ {
		count[nums[i]]++
	}

	index := 0
	for i := 0; i < len(count); i++ {
		for k := 0; k < count[i]; k++ {
			nums[index] = i
			index++
		}
	}
}
