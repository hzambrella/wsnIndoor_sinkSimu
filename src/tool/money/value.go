// 用于项目当中处理各种金额关系。例如：精度转换、汇率转换等.
//
// 参考资料：
// TODO:
//	提供汇率转换。
// BUG(New):
//	未实现四舍五入算法。
//  浮点数运算会有一定的误差
package money

import (
	"database/sql"
	"fmt"
	"math"
	"strconv"
	"errors"

)

const (
	MaxValue             = Value(math.MaxFloat64)             // 金额的最大值
	SmallestNonzeroValue = Value(math.SmallestNonzeroFloat64) // 金额最小可以表示的浮点数。
)

// 金额的统一实体单位。
//
type Value float64

// 创建一个Value的实体单位
//
// 注意：
//	实体单位输入值将依据最小的浮点进行转换。
//	转换时使用四舍五入的算法进行。
// 参数：
//	amount -- 金额原值。
// 返回：
//	返回一个金额值
// BUG：
//	未进行四舍五入转换
//
func New(amount float64) Value {
	return Value(amount)
}

var (
	ErrInput = errors.New("error input")
)

// 转换字符金额到Value，如果转化失败，返回失败的错误
func Parse(amount string) (Value, error) {
	f, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0, ErrInput
	}
	return New(f), nil
}

func (v *Value) Scan(i interface{}) error {
	val := &sql.NullFloat64{}
	if err := val.Scan(i); err != nil {
		return err
	}
	*v = Value(val.Float64)
	return nil
}

// 对两个金额进行比较
// 仅对最多四位精度进行比较
//
// 参数说明：
//	other -- 需要比较的另外一个金额值。
// 返回：
//	0 -- 两个值相等。
//	1 -- 所输入的值小于原值。
//	-1 -- 所输入的值大于原值。
//
func (v Value) Cmp(other Value) int {
	// 仅对最多四位精度进行比较
	return Cmp(float64(v), float64(other), 4)
}

// 取负值
func (v Value) Minus() Value {
	return Value(-v)
}

// 将另一个金额值增加到原值当中。
//
// 注意：
//	增加操作会改变原值的大小。
// 参数：
// 	other -- 需要增加的金额值。
// 返回：
//	返回变更的值。
//
func (v *Value) Add(other Value) *Value {
	*v += other
	return v
}

// 将原值当中减去一个值。
//
// 注意：
//	增加操作会改变原值的大小。
// 参数：
// 	other -- 需要增加的金额值。
// 返回：
//	返回变更值。
//
func (v *Value) Sub(other Value) *Value {
	*v -= other
	return v
}

// 浮点数格式化，参数为精度值
// 使用四舍五入的算法进行格式化输出
func (v Value) Format(prec uint) string {
	return fmt.Sprintf("%."+fmt.Sprintf("%d", prec)+"f", Round(float64(v), prec))
}

// 输出金额值的字符串类型。
// 输出的位数由实际值决定
func (v Value) String() string {
	return fmt.Sprint(Round(float64(v), 4))
}

// 输出金额值的浮点值类型余额。
func (v Value) Amount() float64 {
	return float64(v)
}
