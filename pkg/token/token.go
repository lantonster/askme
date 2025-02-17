package token

import "github.com/google/uuid"

// GenerateToken 函数用于生成一个新的 UUID 作为令牌。
//
// 返回: 生成的 UUID 字符串
func GenerateToken() string {
	// 生成一个版本 7 的 UUID
	uid, _ := uuid.NewV7()
	// 返回 UUID 的字符串表示
	return uid.String()
}
