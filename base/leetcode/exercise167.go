package main

import "fmt"

/**
在有序数组中找到两数之和
https://leetcode-cn.com/problems/two-sum-ii-input-array-is-sorted/
解法：使用对撞指针法
*/
func main() {
	nums := []int{2, 7, 11, 15}
	result := twoSum(nums, 9)
	fmt.Printf("%v", result)
}

func twoSum(numbers []int, target int) []int {
	result := make([]int, 2)
	left := 0
	right := len(numbers) - 1
	for {
		if left >= right {
			break
		}
		sum := numbers[left] + numbers[right]
		if sum == target {
			result[0] = left + 1
			result[1] = right + 1
			break
		} else if sum < target {
			left++
		} else {
			right--
		}
	}
	return result
}
