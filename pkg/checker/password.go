package checker

import (
	"errors"
	"regexp"
	"strings"

	"github.com/lantonster/askme/pkg/reason"
)

const (
	PasswordStrengthLow = iota + 1
	PasswordStrengthMid
	PasswordStrengthHig
)

// CheckPassword 函数用于检查密码是否符合特定规则和强度要求。
//
// 参数:
//   - password: 要检查的密码字符串
//   - strength: 要求的密码强度等级
//
// 返回: 可能返回表示密码不符合要求的错误
func CheckPassword(password string, strength int) error {
	// 检查密码是否包含空格，如果包含则返回错误
	if strings.Contains(password, " ") {
		return errors.New(reason.PasswordCannotContainSpaces)
	}

	level := 0
	patternList := []string{`[0-9]+`, `[a-z]+`, `[A-Z]+`, `[~!@#$%^&*?_-]+`}
	// 遍历正则表达式模式列表，检查密码是否匹配
	for _, pattern := range patternList {
		match, _ := regexp.MatchString(pattern, password)
		if match {
			level++
		}
	}

	// 如果密码的强度等级低于要求，则返回错误
	if level < strength {
		return errors.New(reason.PasswordStrengthTooLow)
	}
	return nil
}
