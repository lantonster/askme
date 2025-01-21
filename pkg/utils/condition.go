package utils

// Ternary 是一个三元操作符函数，根据条件返回两个可能的值中的一个。
// 这个函数的存在是为了提供一个与传统编程语言中的三元操作符类似的语法结构。
//
// 参数:
//   - condition - 一个布尔值，用于决定返回哪个值。
//   - ifOutput - 如果条件为真，则返回的值。
//   - elseOutput - 如果条件为假，则返回的值。
//
// 返回值:
//   - 返回类型为 T，是 ifOutput 和 elseOutput 的共通类型。
//
// 函数根据条件返回 ifOutput 或 elseOutput。
func Ternary[T any](condition bool, ifOutput T, elseOutput T) T {
	// 如果条件为真，则返回 ifOutput。
	if condition {
		return ifOutput
	}

	// 如果条件为假，则返回 elseOutput。
	return elseOutput
}
