package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lantonster/askme/internal/constant"
	"github.com/lantonster/askme/internal/model"
	"github.com/lantonster/askme/internal/schema"
	"github.com/lantonster/askme/internal/service"
	"github.com/lantonster/askme/pkg/errors"
	"github.com/lantonster/askme/pkg/handler"
	"github.com/lantonster/askme/pkg/log"
	"github.com/lantonster/askme/pkg/reason"
)

var ctxUUIDKey = "ctxUuidKey"

type AuthMiddleware struct {
	*service.Service
}

func NewAuthMiddleware(service *service.Service) *AuthMiddleware {
	return &AuthMiddleware{Service: service}
}

// EjectUserBySiteInfo 函数根据站点登录信息和用户信息进行用户权限判断
func (mid *AuthMiddleware) EjectUserBySiteInfo(c *gin.Context) {
	// 获取站点登录信息
	login, err := mid.SiteInfoService().GetSiteLogin(c)
	if err != nil {
		log.WithContext(c).Errorf("获取站点登录信息失败: %v", err)
		c.Next()
		return
	}

	// 如果站点不需要登录，直接跳过
	if !login.LoginRequired {
		c.Next()
		return
	}

	// 尝试从 Cookie 中获取用户信息
	user := GetUserInfoFromContext(c)
	if user == nil {
		handler.Response(c, errors.Unauthorized(reason.UnauthorizedError), nil)
		c.Abort()
		return
	}

	// 用户邮箱未通过验证，拒绝访问
	if user.EmailStatus != model.EmailStatusAvailable {
		handler.Response(c, errors.Forbidden(reason.EmailNeedToBeVerified), &schema.ForbiddenRes{Type: schema.ForbiddenReasonTypeInactive})
		c.Abort()
		return
	}

	c.Next()
}

// NoAuth 认证中间件，如果用户登陆，将用户信息设置到上下文中；如果用户未登陆，继续处理后续中间件或路由。
func (mid *AuthMiddleware) NoAuth(c *gin.Context) {
	// 提取令牌，若令牌为空，继续处理后续中间件或路由
	token := ExtractToken(c)
	if len(token) == 0 {
		c.Next()
		return
	}

	// 获取用户缓存信息，若获取用户信息出错，继续处理后续中间件或路由
	user, err := mid.AuthService().GetUserCacheInfo(c, token)
	if err != nil {
		c.Next()
		return
	}

	// 将用户信息设置到上下文中``
	c.Set(ctxUUIDKey, user)

	// 继续处理后续中间件或路由
	c.Next()
}

// MustAuthAndAccountAvailable 严格的认证中间件，要求令牌有效、账户可用且用户状态正常
// 参数:
//   - c: Gin 上下文
func (mid *AuthMiddleware) MustAuthAndAccountAvailable(c *gin.Context) {
	// 提取令牌，若令牌为空，响应未授权错误并终止请求
	token := ExtractToken(c)
	if len(token) == 0 {
		handler.Response(c, errors.Unauthorized(reason.UnauthorizedError), nil)
		c.Abort()
		return
	}

	// 获取用户缓存信息，若获取出错，响应错误并终止请求
	user, err := mid.AuthService().GetUserCacheInfo(c, token)
	if err != nil {
		handler.Response(c, err, nil)
		c.Abort()
		return
	}

	// 若用户邮箱未通过验证，响应禁止访问错误并终止请求
	if user.EmailStatus != model.EmailStatusAvailable {
		handler.Response(c, errors.Forbidden(reason.EmailNeedToBeVerified), &schema.ForbiddenRes{Type: schema.ForbiddenReasonTypeInactive})
		c.Abort()
		return
	}

	// 若用户状态为被暂停，响应禁止访问错误
	if user.UserStatus == model.UserStatusSuspended {
		handler.Response(c, errors.Forbidden(reason.UserSuspended), &schema.ForbiddenRes{Type: schema.ForbiddenReasonTypeUserSuspended})
	}

	// 若用户状态为已删除，响应未授权错误并终止请求
	if user.UserStatus == model.UserStatusDeleted {
		handler.Response(c, errors.Unauthorized(reason.UnauthorizedError), nil)
		c.Abort()
		return
	}

	// 将用户信息设置到上下文中
	c.Set(ctxUUIDKey, user)

	// 继续处理后续中间件或路由
	c.Next()
}

// MustAuthWithoutAccountAvailable 认证中间件，要求用户必须登录但不要求账户可用。
func (mid *AuthMiddleware) MustAuthWithoutAccountAvailable(c *gin.Context) {
	// 提取令牌，若令牌为空，响应未授权错误并终止请求处理
	token := ExtractToken(c)
	if len(token) == 0 {
		handler.Response(c, errors.Unauthorized(reason.UnauthorizedError), nil)
		c.Abort()
		return
	}

	// 获取用户缓存信息，若获取用户信息出错，响应错误并终止请求处理
	user, err := mid.AuthService().GetUserCacheInfo(c, token)
	if err != nil {
		handler.Response(c, err, nil)
		c.Abort()
		return
	}

	// 若用户状态为已删除，响应未授权错误并终止请求处理
	if user.UserStatus == model.UserStatusDeleted {
		log.WithContext(c).Infof("用户 [%d] 已被删除", user.UserId)
		handler.Response(c, errors.Unauthorized(reason.UnauthorizedError), nil)
		c.Abort()
		return
	}

	// 将用户信息设置到上下文中
	c.Set(ctxUUIDKey, user)

	// 继续处理后续的中间件或路由处理函数
	c.Next()
}

// VisitAuth  认证中间件，要求必须有可用的访问令牌。
func (mid *AuthMiddleware) VisitAuth(c *gin.Context) {
	// 对于特定路径，直接跳过认证，继续处理后续中间件或路由
	if strings.HasPrefix(c.Request.URL.Path, "/uploads/branding/") {
		c.Next()
		return
	}

	// 获取站点登录信息
	login, err := mid.SiteInfoService().GetSiteLogin(c)
	if err != nil {
		return
	} else if !login.LoginRequired {
		c.Next()
		return
	}

	// 尝试从 Cookie 中获取用户访问令牌
	token, err := c.Cookie(constant.CookieKeyUserVisitToken)
	if err != nil || len(token) == 0 {
		// 若获取失败或令牌为空，终止请求并重定向到 403 页面
		c.Abort()
		c.Redirect(http.StatusFound, "/403")
		return
	}

	// 检查访问令牌是否有效，若无效则终止请求并重定向到 403 页面
	if !mid.AuthService().CheckVisitToken(c, token) {
		c.Abort()
		c.Redirect(http.StatusFound, "/403")
		return
	}

	c.Next()
}

// GetUserInfoFromContext 函数从 Gin 上下文获取用户信息。
func GetUserInfoFromContext(c *gin.Context) *model.UserInfo {
	userInfo, exist := c.Get(ctxUUIDKey)
	if !exist {
		return nil
	}

	u, ok := userInfo.(*model.UserInfo)
	if !ok {
		return nil
	}

	return u
}

// GetUserIsAdminModerator 检查 Gin 上下文中的用户是否为管理员或版主。
func GetUserIsAdminModerator(c *gin.Context) bool {
	// 从上下文中获取用户信息，如果用户信息不存在，返回 false
	userInfo, exist := c.Get(ctxUUIDKey)
	if !exist {
		return false
	}

	// 将获取到的用户信息转换为指定的用户信息类型，如果转换不成功，返回 false
	user, ok := userInfo.(*model.UserInfo)
	if !ok {
		return false
	}

	// 如果用户的角色 ID 是管理员或版主的 ID，返回 true
	if user.RoleId == model.RoleIdAdminID || user.RoleId == model.RoleIdModeratorID {
		return true
	}

	return false
}

// ExtractToken 函数从 Gin 上下文的请求头或查询参数中提取令牌。
func ExtractToken(c *gin.Context) (token string) {
	token = c.GetHeader("Authorization")
	if len(token) == 0 {
		token = c.Query("Authorization")
	}
	return strings.TrimPrefix(token, "Bearer ")
}
