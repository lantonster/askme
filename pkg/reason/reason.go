package reason

const (
	// 成功
	Success = "base.success"

	// 未知错误
	UnknownError = "base.unknown"

	// 参数错误
	RequestFormatError = "base.request_format_error"

	// 未登录
	UnauthorizedError = "base.unauthorized_error"

	// 数据库错误
	DatabaseError = "base.database_error"

	// 禁止访问
	ForbiddenError = "base.forbidden_error"

	// 重复请求
	DuplicateRequestError = "base.duplicate_request_error"
)

const (
	// 密码不允许包含空格
	PasswordCannotContainSpaces = "error.password.space_invalid"

	// 密码强度过低
	PasswordStrengthTooLow = "error.password.strength_too_low"

	// site_info 不存在
	SiteInfoNotFound = "error.site_info.not_found"

	// 不允许用户注册
	UserRegistrationNotAllowed = "error.user.registration_not_allowed"
)
