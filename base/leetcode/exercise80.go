package main

import "fmt"

/**
删除排序数组中的重复项II
https://leetcode-cn.com/problems/remove-duplicates-from-sorted-array-ii/
解题关键：只能有两个元素重复，则nums[2] != nums[0]
*/
func main() {
	nums := []int{0, 0, 1, 1, 1, 1, 2, 3, 3}
	fmt.Println(removeDuplicates2(nums))
	fmt.Printf("%v", nums)
}

func removeDuplicates2(nums []int) int {
	if len(nums) < 2 {
		return len(nums)
	}
	k := 2
	for i := 2; i < len(nums); i++ {
		if nums[i] != nums[k-2] {
			nums[k] = nums[i]
			k++
		}
	}
	return k
}
