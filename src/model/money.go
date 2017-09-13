package model

import (
	"fmt"
	"strconv"
)

// Money 钱钱钱
type Money float32

func (m *Money) add(m1 float32) {
	*m += Money(m1)
}

func (m *Money) multiply(m1 float32) {
	*m = *m * Money(m1)
}

// 银行家四舍五入算法
func (m *Money) getFen() int64 {
	f, _ := strconv.ParseFloat(fmt.Sprintf("%0.2f", *m), 32)
	return int64(f * 100)
}
