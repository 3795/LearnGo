package main

import (
	"fmt"
	"math"
)

/**
盛最多水的容器
https://leetcode-cn.com/problems/container-with-most-water/
解析：使用对撞指针法，两个指针向中间移动，移动过程中，舍弃高度较短的那边，因为两个指针向中间移动，间距在减少，
想要面积最大化，则必须选择高度最高的元素
*/

func main() {
	height := []int{1, 8, 6, 2, 5, 4, 8, 3, 7}
	fmt.Println(maxArea(height))
}

func maxArea(height []int) int {
	left := 0
	right := len(height) - 1
	maxarea := 0
	for {
		if left > right {
			break
		}
		tmpArea := int(math.Min(float64(height[left]), float64(height[right]))) * (right - left)
		maxarea = int(math.Max(float64(maxarea), float64(tmpArea)))
		if height[left] < height[right] {
			left++
		} else {
			right--
		}
	}
	return maxarea
}
