package model

// VerificationEmail [Redis] 电子邮件验证码相关的内容缓存
type VerificationEmail struct {
	UserId                   int64  `json:"user_id"`                     // 用户 ID
	Email                    string `json:"e_mail"`                      // 电子邮件地址
	BindingKey               string `json:"binding_key"`                 // 绑定键
	SkipValidationLatestCode bool   `json:"skip_validation_latest_code"` // 跳过验证最新的验证码
}
