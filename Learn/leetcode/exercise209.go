package main

import (
	"fmt"
	"math"
)

/**
长度最小子数组
https://leetcode-cn.com/problems/minimum-size-subarray-sum/
解法：使用滑动窗口的方法
*/

func main() {
	nums := []int{2, 3, 1, 2, 4, 3}
	fmt.Println(minSubArrayLen(7, nums))
}

func minSubArrayLen(s int, nums []int) int {
	left := 0
	right := -1
	res := len(nums) + 1
	sum := 0

	for {
		if left >= len(nums) {
			break
		}

		if right+1 < len(nums) && sum < s {
			right++
			sum += nums[right]
		} else {
			sum -= nums[left]
			left++
		}

		if sum >= s {
			res = int(math.Min(float64(res), float64(right-left+1)))
		}
	}

	if res == len(nums)+1 {
		return 0
	}
	return res
}
