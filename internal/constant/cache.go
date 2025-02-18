package constant

import "time"

const (
	// 配置: (string)key -> (*model.Config)config
	CacheKeyConfig  = "askme:config:key:%s"
	CacheTimeConfig = time.Hour

	// 站点信息: (model.SiteInfoType)type -> (*model.SiteInfo)site info
	CacheKeySiteInfo  = "askme:site_info:type:%s"
	CacheTimeSiteInfo = time.Hour

	// 用户信息: (string)access token -> (*model.UserInfo)user info
	CacheKeyUserInfo  = "askme:user:info:access_token:%s"
	CacheTimeUserInfo = 7 * 24 * time.Hour

	// 用户 access token 映射表: (int64)user id -> (map[string]bool)map[access_token]bool
	CacheKeyUserAccessTokenMapping  = "askme:user:access_token:mapping:user_id:%d"
	CacheTimeUserAccessTokenMapping = 7 * 24 * time.Hour

	// 用户 access token: (string)visit token -> (string)access token
	CacheKeyUserAccessToken  = "askme:user:access_token:visit_token:%s"
	CacheTimeUserAccessToken = 7 * 24 * time.Hour

	// 用户最新邮件验证码: (int64)user id -> (string)code
	CacheKeyVerificationEmailLatestCode = "askme:verification_email:latest_code:user_id:%d"

	// 验证邮件: (string)code -> (*model.VerificationEmail)email content
	CacheKeyVerificationEmail  = "askme:verification_email:code:%s"
	CacheTimeVerificationEmail = 10 * time.Minute
)
