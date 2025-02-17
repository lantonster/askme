package random

import (
	"crypto/rand"
	"encoding/hex"
)

// UsernameSuffix 函数用于生成一个随机的用户名后缀。
//
// 返回: 一个由随机字节生成并编码为十六进制字符串的后缀
func UsernameSuffix() string {
	// 创建一个长度为 2 的字节切片用于存储随机字节
	bytes := make([]byte, 2)
	// 从随机源读取随机字节到切片中，忽略可能的错误
	_, _ = rand.Read(bytes)
	// 将随机字节编码为十六进制字符串并返回
	return hex.EncodeToString(bytes)
}
