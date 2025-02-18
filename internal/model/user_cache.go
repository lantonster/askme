package model

// UserInfo [Redis] 用户缓存信息
type UserInfo struct {
	UserId      int64  `json:"user_id"`
	RoleId      int64  `json:"role_id"`
	UserStatus  string `json:"user_status"`
	EmailStatus string `json:"email_status"`
	ExternalId  string `json:"external_id"`
	VisitToken  string `json:"visit_token"`
}
