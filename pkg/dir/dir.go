package dir

import "os"

// CreateDirIfNotExist 创建一个目录，如果该目录不存在。
//
// 参数:
//   - path: 要创建的目录的路径
//
// 返回: 可能返回创建目录时发生的错误
func CreateDirIfNotExist(path string) error {
	// 使用 os.MkdirAll 函数创建目录及其所有必要的父目录，并设置权限为 os.ModePerm
	return os.MkdirAll(path, os.ModePerm)
}

// CheckFileExist 检查指定路径的文件是否存在。
//
// 参数:
//   - path: 要检查的文件的路径
//
// 返回: 如果文件存在且不是目录则返回 true，否则返回 false
func CheckFileExist(path string) bool {
	// 获取文件或目录的信息
	f, err := os.Stat(path)
	// 如果获取信息没有错误且不是目录，返回 true
	return err == nil && !f.IsDir()
}
