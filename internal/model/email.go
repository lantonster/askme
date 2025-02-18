package model

// VerificationEmail [Redis] 电子邮件验证码相关的内容缓存
type VerificationEmail struct {
	UserID int64  `json:"user_id"` // 用户 ID
	Email  string `json:"e_mail"`  // 电子邮件地址
}
