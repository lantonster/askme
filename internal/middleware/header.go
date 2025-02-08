package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// SetStaticCacheHeader 设置静态文件缓存
func SetStaticCacheHeader(c *gin.Context) {
	if strings.HasPrefix(c.Request.RequestURI, "/static/") {
		c.Header("Cache-Control", "public, max-age=31536000")
	}
}
