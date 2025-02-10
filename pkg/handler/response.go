package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lantonster/askme/pkg/errors"
	"github.com/lantonster/askme/pkg/log"
	"github.com/lantonster/askme/pkg/reason"
)

// ResponseBody response body.
type ResponseBody struct {
	// http code
	Code int `json:"code"`
	// reason key
	Reason string `json:"reason"`
	// response message
	Message string `json:"msg"`
	// response data
	Data any `json:"data"`
}

func NewResponseBody(code int, reason string) *ResponseBody {
	return &ResponseBody{
		Code:   code,
		Reason: reason,
	}
}

func NewresponseBodyError(err *errors.Error, data any) *ResponseBody {
	return &ResponseBody{
		Code:    err.Code,
		Reason:  err.Reason,
		Message: err.Message,
		Data:    data,
	}
}

func NewResponseBodyData(code int, reason string, data any) *ResponseBody {
	return &ResponseBody{
		Code:   code,
		Reason: reason,
		Data:   data,
	}
}

// Response 函数用于处理 Gin 框架中的 HTTP 响应
// 参数:
//   - c: Gin 上下文
//   - err: 可能的错误
//   - data: 要返回的数据
func Response(c *gin.Context, err error, data any) {
	// 如果没有错误，返回成功状态码和包含数据的响应体，并终止后续处理
	if err == nil {
		c.JSON(http.StatusOK, NewResponseBodyData(http.StatusOK, reason.Success, data))
		c.Abort()
		return
	}

	var myErr *errors.Error
	// 未知错误
	if !errors.As(err, &myErr) {
		log.WithContext(c).Errorf("http 响应未知错误: %v", err)
		c.JSON(http.StatusInternalServerError, NewResponseBody(http.StatusInternalServerError, reason.UnknownError))
		c.Abort()
		return
	}

	// 服务器内部错误
	if errors.IsInternalServer(myErr) {
		log.WithContext(c).Errorf("http 响应服务器内部错误: %v", err)
	}

	// 创建并返回包含错误和数据的响应体
	body := NewresponseBodyError(myErr, data)
	c.JSON(myErr.Code, body)
	c.Abort()
}
