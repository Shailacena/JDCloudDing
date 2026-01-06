package util

import "github.com/shopspring/decimal"

// ToDecimal 价格转换为2位小数
func ToDecimal(num float64) float64 {
	amount, _ := decimal.NewFromFloat(num).Round(2).Float64()

	return amount
}

// F32ToDecimal 价格转换为2位小数
func F32ToDecimal(num float32) float64 {
	amount, _ := decimal.NewFromFloat32(num).Round(2).Float64()

	return amount
}