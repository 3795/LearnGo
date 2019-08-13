package main

import "fmt"

// https://leetcode-cn.com/problems/single-number/
func main() {
	nums := []int{4, 1, 2, 1, 2}
	fmt.Println(singleNumber(nums))
}

func singleNumber(nums []int) int {
	m := make(map[int]int)
	for _, value := range nums {
		_, ok := m[value]
		if !ok {
			m[value] = 1
		} else {
			delete(m, value)
		}
	}

	for key := range m {
		if m[key] == 1 {
			return key
		}
	}
	return 0
}
