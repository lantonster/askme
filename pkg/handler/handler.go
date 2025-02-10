package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lantonster/askme/pkg/errors"
	"github.com/lantonster/askme/pkg/log"
	"github.com/lantonster/askme/pkg/reason"
	"github.com/lantonster/askme/pkg/validator"
)

// BindAndCheck 函数用于在 Gin 框架中绑定请求参数并进行校验。
//
// 参数:
//   - c: Gin 上下文
//   - req: 要绑定和校验的请求结构体
//
// 返回:
//   - bool: 如果绑定或校验失败则为 true，否则为 false
func BindAndCheck(c *gin.Context, req any) (abort bool) {
	// 尝试绑定请求参数到给定的结构体
	if err := c.ShouldBind(req); err != nil {
		log.WithContext(c).Errorf("http 解析请求参数失败: %v", err)
		Response(c, errors.New(http.StatusBadRequest, reason.RequestFormatError), err.Error())
		return true
	}

	// 进行请求参数的校验
	if fields, err := validator.Check(c, req); err != nil {
		log.WithContext(c).Errorf("http 校验请求参数失败: %v", err)
		Response(c, err, fields)
		return true
	}

	return false
}
