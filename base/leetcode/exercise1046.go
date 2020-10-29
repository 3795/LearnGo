package main

import (
	"fmt"
	"sort"
)

/**
最后一块石头的重量
https://leetcode-cn.com/problems/last-stone-weight/submissions/
*/

func main() {
	stones := []int{1, 1, 3, 4}
	result := lastStoneWeight(stones)
	fmt.Println(result)
}

func lastStoneWeight(stones []int) int {
	weight := 0
	length := len(stones)
	for i := 0; i < length-1; i++ {
		sort.Ints(stones)
		weight = stones[length-1] - stones[length-2]
		stones[length-1] = weight
		stones[length-2] = 0
	}
	return stones[length-1]
}
