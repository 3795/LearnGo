package main

import "fmt"

/**
找到字符串中所有字母的异位词
https://leetcode-cn.com/problems/find-all-anagrams-in-a-string/
*/

func main() {
	//s := "cbaebabacd"
	p := "abc"
	fmt.Println(string(p[0]))
}

func findAnagrams(s string, p string) []int {
	index := 0
	var res []int
	for {
		if index > len(s)-len(p) {
			break
		}

		// todo 完成这道题
	}

	return res
}
