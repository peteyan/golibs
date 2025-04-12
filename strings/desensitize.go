package strings

import (
	"regexp"
	"strings"
)

// DesensitizeCommon 脱敏字符串，如果字符串大于6个字符，则只显示前三个字符和后三个字符，其他字符用星号代替；否则返回星号
func DesensitizeCommon(source string) string {
	if len(source) > 6 {
		return source[:3] + strings.Repeat("*", len(source)-6) + source[len(source)-3:]
	}
	return strings.Repeat("*", len(source))
}

// DesensitizeEmail 脱敏电子邮件地址
func DesensitizeEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email // 非法的邮箱格式，直接返回
	}
	username := parts[0]
	if len(username) > 2 {
		username = username[:2] + strings.Repeat("*", len(username)-2)
	}
	return username + "@" + parts[1]
}

// DesensitizeCreditCard 脱敏信用卡号码
func DesensitizeCreditCard(cardNumber string) string {
	if len(cardNumber) < 12 {
		return cardNumber // 非法的信用卡号码，直接返回
	}
	return strings.Repeat("*", len(cardNumber)-4) + cardNumber[len(cardNumber)-4:]
}

// DesensitizePhoneNumber 脱敏手机号码
func DesensitizePhoneNumber(phone string) string {
	if len(phone) < 11 {
		return phone // 非法的手机号码，直接返回
	}
	return phone[:3] + strings.Repeat("*", len(phone)-7) + phone[len(phone)-4:]
}

// DesensitizeName 脱敏姓名
func DesensitizeName(name string) string {
	if len(name) <= 1 {
		return name // 如果姓名只有一个字，直接返回
	}
	return string(name[0]) + strings.Repeat("*", len(name)-1)
}

// DesensitizeCustom 使用正则表达式进行自定义脱敏
func DesensitizeCustom(input, pattern, replacement string) string {
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(input, replacement)
}
