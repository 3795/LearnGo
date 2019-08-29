package array_util

import (
	"math/rand"
	"time"
)

/**
生成随机数数组
*/
func GenArr(n int) []int {
	arr := make([]int, n)
	rand.Seed(time.Now().Unix())
	for i := 0; i < n; i++ {
		arr[i] = rand.Intn(1000)
	}
	return arr
}
