package gravatar

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// GetAvatarURL 根据给定的基础 URL 和电子邮件地址生成头像的 URL。
//
// 参数:
//   - baseURL: 头像的基础 URL
//   - email: 用于生成哈希值的电子邮件地址
//
// 返回: 生成的头像 URL
func GetAvatarURL(baseURL, email string) string {
	// 对电子邮件地址进行处理并计算 SHA-256 哈希值
	hasher := sha256.Sum256([]byte(strings.TrimSpace(email)))
	// 将哈希值编码为十六进制字符串
	hash := hex.EncodeToString(hasher[:])
	// 组合基础 URL 和哈希值得到最终的头像 URL
	return baseURL + hash
}
