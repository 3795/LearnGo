package main

import "fmt"

/**
合并两个有序数组
https://leetcode-cn.com/problems/merge-sorted-array/
*/

func main() {
	nums1 := []int{1, 2, 3, 0, 0, 0}
	nums2 := []int{2, 5, 6}
	result := *merge(nums1, 3, nums2, 3)
	for _, i := range result {
		fmt.Println(i)
	}
}

func merge(nums1 []int, m int, nums2 []int, n int) *[]int {
	tmp := make([]int, m)
	copy(tmp, nums1[:m])
	p := 0
	q := 0
	i := 0
	for {
		if p >= m || q >= n {
			break
		}
		if tmp[p] <= nums2[q] {
			nums1[i] = tmp[p]
			p++
		} else {
			nums1[i] = nums2[q]
			q++
		}
		i++
	}

	if p < m {
		for i := p; i < m; i++ {
			nums1[i+n] = tmp[i]
		}
	}

	if q < n {
		for i := q; i < n; i++ {
			nums1[i+p] = nums2[i]
		}
	}

	return &nums1
}
