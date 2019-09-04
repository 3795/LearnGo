package main

import "fmt"

/**
两个数组的交集
https://leetcode-cn.com/problems/intersection-of-two-arrays-ii/
解析：使用计数查找的方法
*/

func main() {
	nums1 := []int{4, 9, 5}
	nums2 := []int{9, 4, 9, 8, 4}
	res := intersect(nums1, nums2)
	fmt.Printf("%v", res)
}

func intersect(nums1 []int, nums2 []int) []int {
	m := make(map[int]int)
	var res []int
	for _, n := range nums2 {
		value, ok := m[n]
		if ok {
			m[n] = value + 1
		} else {
			m[n] = 1
		}
	}

	for _, item := range nums1 {
		v, ok := m[item]
		if ok && v > 0 {
			res = append(res, item)
			m[item] = v - 1
		}
	}

	return res
}
