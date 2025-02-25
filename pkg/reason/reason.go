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
	// 配置不存在
	ConfigNotFound = "error.config.not_found"

	// 邮箱需要验证
	EmailNeedToBeVerified = "error.email.need_to_be_verified"

	// 邮箱域名不合法
	EmailIllegalDomainError = "error.email.illegal_email_domain_error"

	// 邮箱或密码错误
	EmailOrPasswordWrong = "error.object.email_or_password_incorrect"

	// 邮箱验证链接已过期
	EmailVerifyURLExpired = "error.email.verify_url_expired"

	// 不允许通过密码登录
	NotAllowedLoginViaPassword = "error.user.not_allowed_login_via_password"

	// 密码不允许包含空格
	PasswordCannotContainSpaces = "error.password.space_invalid"

	// 密码强度过低
	PasswordStrengthTooLow = "error.password.strength_too_low"

	// site_info 不存在
	SiteInfoNotFound = "error.site_info.not_found"

	// 邮箱已被注册
	UserEmailDuplicate = "error.email.duplicate"

	// 用户名已被注册
	UsernameDuplicate = "error.user.username_duplicate"

	// 用户名不合法
	UsernameInvalid = "error.user.username_invalid"

	// 用户不存在
	UserNotFound = "error.user.not_found"

	// 不允许用户注册
	UserRegistrationNotAllowed = "error.user.registration_not_allowed"

	// 用户被禁用
	UserSuspended = "error.user.suspended"
)
