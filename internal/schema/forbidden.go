package schema

const (
	ForbiddenReasonTypeInactive      = "inactive"    // 用户未激活
	ForbiddenReasonTypeURLExpired    = "url_expired" // 链接已过期
	ForbiddenReasonTypeUserSuspended = "suspended"   // 用户被禁用
)

type ForbiddenRes struct {
	Type string `json:"type" enums:"inactive,url_expired,suspended"`
}
