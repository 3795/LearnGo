package main

import "fmt"

/**
求众数
https://leetcode-cn.com/problems/majority-element/
*/

func main() {
	nums := []int{2, 2, 1, 1, 1, 2, 2}
	major := majorityElement(nums)
	fmt.Println(major)
}

func majorityElement(nums []int) int {
	major := nums[0]
	count := 1
	for i := 1; i < len(nums); i++ {
		if major == nums[i] {
			count++
		} else {
			count--
			if count <= 0 {
				major = nums[i+1]
			}
		}
	}
	return major
}
