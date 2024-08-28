package utils

import (
	"math/rand"
	"time"
)

// RangeRandom 返回指定范围内的随机整数。[min,max),不包含max
func RangeRandom(min, max int) (number int) {
	//创建随机种子
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	number = r.Intn(max-min) + min
	return number
}
