package errors

import (
	"bytes"
	"fmt"
	"runtime"
)

// LogStack 返回从 start 到 end 的调用栈信息
func LogStack(start, end int) string {
	stack := bytes.Buffer{}
	for i := start; i < end || end <= 0; i++ {
		pc, str, line, _ := runtime.Caller(i)
		if line == 0 {
			break
		}
		stack.WriteString(fmt.Sprintf("%s:%d %s\n", str, line, runtime.FuncForPC(pc).Name()))
	}
	return stack.String()
}
