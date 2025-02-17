package checker

import (
	"encoding/json"
	"regexp"

	"github.com/lantonster/askme/configs"
)

var (
	// usernameReg 要求用户名由 2 到 30 个小写字母、数字、点、下划线和连字符组成
	usernameReg = regexp.MustCompile(`^[a-z0-9._-]{2,30}$`)

	// reservedUsernameMapping 用于存储保留的用户名映射
	reservedUsernameMapping = make(map[string]bool)
)

// IsInvalidUsername 函数用于检查给定的用户名是否无效。
//
// 参数:
//   - username: 要检查的用户名
//
// 返回: 如果用户名不符合定义的正则表达式格式则返回 true，否则返回 false
func IsInvalidUsername(username string) bool {
	// 通过正则表达式匹配来判断用户名是否有效
	return !usernameReg.MatchString(username)
}

func init() {
	// TODO 读取目标文件

	var usernames []string
	_ = json.Unmarshal(configs.ReservedUsernames, &usernames)
	for _, username := range usernames {
		reservedUsernameMapping[username] = true
	}
}

// IsReservedUsername 检查给定的用户名是否为保留用户名。
//
// 参数:
//   - username: 要检查的用户名
//
// 返回: 如果是保留用户名则返回 true，否则返回 false
func IsReservedUsername(username string) bool {
	return reservedUsernameMapping[username]
}
