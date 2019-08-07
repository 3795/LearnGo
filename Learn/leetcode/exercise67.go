package main

import (
	"fmt"
	"strconv"
)

/**
二进制相加
*/
func addBinary(a string, b string) string {
	var ans string
	ca := 0
	sum := 0
	i := len(a) - 1
	j := len(b) - 1
	for {
		if i < 0 && j < 0 {
			break
		}
		sum = ca
		if i >= 0 {
			t1, _ := strconv.Atoi(string([]rune(a)[i]))
			sum += t1
		} else {
			sum += 0
		}
		if j >= 0 {
			t2, _ := strconv.Atoi(string([]rune(b)[j]))
			sum += t2
		} else {
			sum += 0
		}
		ans += strconv.Itoa(sum % 2)
		ca = sum / 2

		i--
		j--
	}
	if ca == 1 {
		ans += "1"
	}
	return reverse(ans)
}

func reverse(str string) string {
	rs := []rune(str)
	length := len(rs)
	var tt []rune

	for i := length - 1; i >= 0; i-- {
		tt = append(tt, rs[i])
	}
	return string(tt[0:])
}

func main() {
	fmt.Println(addBinary("11", "1"))
}
