package checkinput

import (
	"regexp"
)

//TODO:有问题，不要用！
const(
	checkNumber="-?[0-9]+.*[0-9]*"
)

/**
 * 检验是否是数字
 *
 * @param regex
 * @param input
 * @return
 */
func CheckNumber(input string) bool {
	return match(checkNumber,input)
}

/**
 * 匹配函数
 *
 * @param regex
 * @param input
 * @return
 */
func match(regex string, input string) bool {
	rex := regexp.MustCompile(regex)
	return rex.MatchString(input)
}