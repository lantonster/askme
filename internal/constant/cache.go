package constant

import "time"

const (
	// 配置: (string)key -> (*model.Config)config
	CacheKeyConfig  = "askme:config:key:%s"
	CacheTimeConfig = time.Hour

	// 站点信息: (model.SiteInfoType)type -> (*model.SiteInfo)site info
	CacheKeySiteInfo  = "askme:site_info:type:%s"
	CacheTimeSiteInfo = time.Hour

	// 用户最新邮件验证码: (int64)user id -> (string)code
	CacheKeyVerificationEmailLatestCode = "askme:verification_email:latest_code:user_id:%d"

	// 验证邮件: (string)code -> (*model.VerificationEmail)email content
	CacheKeyVerificationEmail  = "askme:verification_email:code:%s"
	CacheTimeVerificationEmail = 10 * time.Minute
)
