package reason

const (
	// 成功
	Success = "success"

	// 未知错误
	UnknownError = "unknown"

	// 参数错误
	RequestFormatError = "request_format_error"

	// 未登录
	UnauthorizedError = "unauthorized_error"

	// 数据库错误
	DatabaseError = "database_error"

	// 禁止访问
	ForbiddenError = "forbidden_error"

	// 重复请求
	DuplicateRequestError = "duplicate_request_error"
)
