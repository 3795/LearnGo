package main

import (
	"fmt"
	"regexp"
	"strings"
)

/**
验证回文字符串
https://leetcode-cn.com/problems/valid-palindrome/
*/
func main() {
	str := "aaaa........."
	fmt.Println(isPalindrome(str))
}

func isPalindrome(s string) bool {
	if len(s) == 0 {
		return true
	}
	s = strings.ToLower(s)
	r := []rune(s)
	re := regexp.MustCompile(`[a-z0-9]`)
	left := 0
	right := len(s) - 1
	leftStr := ""
	rightStr := ""
	for {
		if left > right {
			break
		}

		leftStr = string(r[left])
		rightStr = string(r[right])
		if !re.MatchString(leftStr) {
			left++
			continue
		}
		if !re.MatchString(rightStr) {
			right--
			continue
		}
		if leftStr == rightStr {
			left++
			right--
		} else {
			return false
		}
	}
	return true
}
