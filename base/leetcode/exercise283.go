package main

import "fmt"

/**
移动零
https://leetcode-cn.com/problems/move-zeroes/
*/
func main() {
	nums := []int{0, 1, 0, 3, 12}
	moveZeroes(nums)
	fmt.Printf("%v", nums)
}

func moveZeroes(nums []int) {
	work := 0
	for i := 0; i < len(nums); i++ {
		if nums[i] != 0 {
			// 与work进行置换
			if i != work {
				swap(nums, i, work)
			}
			work++
		}
	}
}

func swap(nums []int, i, j int) {
	t := nums[i]
	nums[i] = nums[j]
	nums[j] = t
}
