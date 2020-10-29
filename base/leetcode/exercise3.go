package main

import (
	"fmt"
	"math"
)

/**
无重复字符的最长子串
https://leetcode-cn.com/problems/longest-substring-without-repeating-characters/
解法：使用滑动窗口
*/
func main() {
	s := "bbbbbb"
	fmt.Println(lengthOfLongestSubstring(s))
}

func lengthOfLongestSubstring(s string) int {
	freq := make([]int, 256) // 设置大小为256的原因是a-z的ASCII值为97-122，用不完的地方就空着
	left := 0
	right := -1
	res := 0
	// 每次在循环中维护freq，并记录窗口中是否能找到一个新的最优值
	for {
		if left >= len(s) {
			break
		}

		if right+1 < len(s) && freq[s[right+1]] == 0 {
			right += 1
			freq[s[right]]++
		} else {
			freq[s[left]]--
			left += 1
		}

		res = int(math.Max(float64(res), float64(right-left+1)))
	}
	return res
}
