package main

import "fmt"

func main() {
	strs := []byte{'h', 'e', 'l', 'l', 'o'}
	reverseString(strs)
	fmt.Printf("%s", strs)
}

func reverseString(s []byte) {
	left := 0
	right := len(s) - 1
	for {
		if left > right {
			break
		}
		swapByte(s, left, right)
		left++
		right--
	}
}

func swapByte(s []byte, i, j int) {
	tmp := s[i]
	s[i] = s[j]
	s[j] = tmp
}
