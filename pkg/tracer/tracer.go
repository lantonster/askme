package tracer

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const TraceIdKey = "askme_trace_id"

// Trace 在 http 请求处理链中插入跟踪 id。
//
// 该方法旨在为每个 HTTP 请求分配一个唯一的跟踪 id，以便于后续的日志跟踪和问题排查。
// 它利用 gin 框架的中间件机制，在请求处理的早期阶段注入跟踪 id，确保该 id 在整个请求处理过程中可用。
func Trace(c *gin.Context) {
	// 从请求头中尝试获取跟踪 ID。
	traceId := c.GetHeader(TraceIdKey)

	// 如果请求头中没有跟踪 ID，则生成一个新的跟踪 ID。
	if traceId == "" {
		traceId = uuid.New().String()
	}

	// 使用生成或获取的跟踪 ID 设置日志上下文，以便后续的日志输出中包含该跟踪 ID。
	c.Set(TraceIdKey, traceId)

	// 将跟踪 ID 添加到响应头中，以便客户端可以获取跟踪 ID。
	c.Writer.Header().Set(TraceIdKey, traceId)

	// 继续处理请求，将控制权交给下一个中间件或处理函数。
	c.Next()
}
