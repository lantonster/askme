package constant

import "time"

const (
	// 用户 visit token: *gin.Context -> (string)visit token
	CookieKeyUserVisitToken  = "visit_token"
	CookieTimeUserVisitToken = 7 * 24 * time.Hour
)
