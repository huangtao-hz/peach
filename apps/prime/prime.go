package main

import (
	"math"
	"peach/utils"
	"time"
)

// squre 求平方
func squre[T ~int](i T) T {
	return i * i
}

// prime 求指定数字内所有素数
func prime(x int) (result []int) {
	var i, j, k int
	// 预估 x 以内素数的个数
	cap := int(float64(x) * 1.2 / math.Log(float64(x)))
	// 分配素数序列的内存
	result = make([]int, 0, cap)
	for i = 2; i < min(squre(2), x+1); i++ {
		result = append(result, i)
	}
	var isPrime bool
	for j = 1; i < x; j++ {
		for ; i < min(x, squre(result[j])); i++ {
			isPrime = true
			for _, k = range result[:j] {
				if i%k == 0 {
					isPrime = false
					break
				}
			}
			if isPrime {
				result = append(result, i)
			}
		}
	}
	return
}

// main 主函数
func main() {
	defer utils.TimeIt(time.Now())
	x := 3000000
	xl := prime(x)
	utils.Printf("count:%,d\n", len(xl))
	//fmt.Println("前100个素数：", xl[:100])
}
