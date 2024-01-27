package common

import "regexp"

// 数据格式验证
func VerifyMobile(phoneNumber string) bool {
	regex := `^1\d{10}$`
	match, err := regexp.MatchString(regex, phoneNumber)
	if err != nil {
		return false
	}

	return match
}

func VerifyEmailFormat(email string) bool {
	//pattern := `\w+([-+.]\w+)@\w+([-.]\w+).\w+([-.]\w+)*` //匹配电子邮箱
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z].){1,4}[a-z]{2,4}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}
