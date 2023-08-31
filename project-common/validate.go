package common

import "regexp"

// 数据格式验证
func VerifyMobile(phoneNumber string) bool {
	pattern := `(13[0-9]|14[0-9]|15[0-9]|16[0-9]|17[0-9]|18[0-9]|19[0-9])`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(phoneNumber)
}
