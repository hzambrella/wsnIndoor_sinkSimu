// 仅适用于浮点数，精度至少为0
package money

import (
	"math"
)

// 比较两人个浮点数
// a -- 第一个浮点数
// b -- 第二个浮点数
// prec -- 比较的最大精度
//
// 返回
// 0 -- 相等
// -1 -- a < b
// 1 -- a > b
func Cmp(a, b float64, prec uint) int {
	c := Round(a, prec) - Round(b, prec)
	//	1 -- 所输入的值小于原值。
	if c > 0 {
		return 1
	}
	//	-1 -- 所输入的值大于原值。
	if c < 0 {
		return -1
	}
	return 0
}

// 四舍五入
func Round(in float64, prec uint) float64 {
	d := math.Pow10(int(prec))
	return math.Floor((in*d)+0.5) / d
}

// 向上舍入
func RoundUp(in float64, prec uint) float64 {
	d := math.Pow10(int(prec))
	return math.Ceil(in*d) / d
}

// 向下舍入
func RoundDown(in float64, prec uint) float64 {
	d := math.Pow10(int(prec))
	return math.Floor(in*d) / d
}
