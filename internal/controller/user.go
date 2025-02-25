package controller

import (
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/lantonster/askme/internal/constant"
	"github.com/lantonster/askme/internal/middleware"
	"github.com/lantonster/askme/internal/schema"
	"github.com/lantonster/askme/internal/service"
	"github.com/lantonster/askme/pkg/errors"
	"github.com/lantonster/askme/pkg/handler"
	"github.com/lantonster/askme/pkg/log"
	"github.com/lantonster/askme/pkg/reason"
)

type UserController struct {
	*service.Service
}

func NewUserController(service *service.Service) *UserController {
	return &UserController{Service: service}
}

// SetVisitCookies 函数用于设置用户访问的 Cookie。
//
// 参数:
//   - c: Gin 上下文
//   - visitToken: 访问令牌
//   - force: 是否强制设置，即使已有 Cookie 且不强制时不会设置
func (ctrl *UserController) SetVisitCookies(c *gin.Context, visitToken string, force bool) {
	// 尝试获取已有的访问令牌 Cookie
	cookie, err := c.Cookie(constant.CookieKeyUserVisitToken)
	if err == nil && len(cookie) > 0 && !force {
		return
	}

	// 获取站点通用信息
	general, err := ctrl.SiteInfoService().GetSiteGeneral(c)
	if err != nil {
		log.WithContext(c).Errorf("获取站点信息失败: %v", err)
		return
	}

	// 解析站点 URL
	url, err := url.Parse(general.SiteUrl)
	if err != nil {
		log.WithContext(c).Errorf("解析站点 URL [%s] 失败: %v", general.SiteUrl, err)
		return
	}

	// 设置新的 Cookie
	c.SetCookie(constant.CookieKeyUserVisitToken, visitToken, int(constant.CookieTimeUserVisitToken), "/", url.Host, true, true)
}

// LoginByEmail godoc
//
//	@Summary		通过邮箱登录
//	@Description	通过邮箱登录
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			body	body		schema.LoginByEmailReq								true	"body"
//	@Success		200		{object}	handler.ResponseBody{data=schema.LoginByEmailRes}	"success"
//	@Router			/askme/api/v1/user/email/login [post]
func (ctrl *UserController) LoginByEmail(c *gin.Context) {
	req := &schema.LoginByEmailReq{}
	if handler.BindAndCheck(c, req) {
		return
	}

	res, err := ctrl.UserService().LoginByEmail(c, req)
	ctrl.SetVisitCookies(c, res.VisitToken, true)
	handler.Response(c, err, res)
}

// RegisterUserByEmail godoc
//
//	@Summary		通过邮箱注册账号
//	@Description	通过邮箱注册账号
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			body	body		schema.RegisterUserByEmailReq								true	"body"
//	@Success		200		{object}	handler.ResponseBody{data=schema.RegisterUserByEmailRes}	"success"
//	@Router			/askme/api/v1/user/email/register [post]
func (ctrl *UserController) RegisterUserByEmail(c *gin.Context) {
	// 获取站点登录信息，如果获取失败则响应错误并返回
	loginInfo, err := ctrl.SiteInfoService().GetSiteLogin(c)
	if err != nil {
		handler.Response(c, err, nil)
		return
	}
	// 如果不允许新用户注册或不允许通过电子邮件注册
	if !loginInfo.AllowNewRegistrations || !loginInfo.AllowEmailRegistrations {
		handler.Response(c, errors.BadRequest(reason.UserRegistrationNotAllowed), nil)
		return
	}

	req := &schema.RegisterUserByEmailReq{}
	// 绑定并检查请求参数，如果绑定和检查出错则返回
	if handler.BindAndCheck(c, req) {
		return
	}
	// 如果登录信息不允许该邮箱注册
	if !loginInfo.IsEmailAllowed(req.Email) {
		handler.Response(c, errors.BadRequest(reason.EmailIllegalDomainError), nil)
		return
	}

	req.IP = c.ClientIP()
	// 调用用户服务进行注册，并根据结果进行响应
	res, fieldErr, err := ctrl.UserService().RegisterUserByEmail(c, req)
	if len(fieldErr) > 0 {
		handler.Response(c, err, fieldErr)
	} else {
		handler.Response(c, err, res)
	}
}

// VerifyEmail godoc
//
//	@Summary		邮箱验证
//	@Description	邮箱验证
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			body	body		schema.VerifyEmailReq								true	"body"
//	@Success		200		{object}	handler.ResponseBody{data=schema.VerifyEmailRes}	"success"
//	@Router			/askme/api/v1/user/email/verification [post]
func (ctrl *UserController) VerifyEmail(c *gin.Context) {
	req := &schema.VerifyEmailReq{}
	if handler.BindAndCheck(c, req) {
		return
	}

	email, res, err := ctrl.EmailService().VerifyUrlExpired(c, req.Code)
	if err != nil {
		handler.Response(c, err, res)
	}

	req.Email = email
	resp, err := ctrl.UserService().VerifyEmail(c, req)
	if err != nil {
		handler.Response(c, err, nil)
		return
	}

	handler.Response(c, nil, resp)
}

// CurrentUserInfo godoc
//
//	@Summary		获取用户信息
//	@Description	获取用户信息
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	handler.ResponseBody{data=schema.GetUserByUserIdRes}	"success"
//	@Router			/askme/api/v1/user/info [get]
func (ctrl *UserController) CurrentUserInfo(c *gin.Context) {
	// 从上下文中获取用户信息
	user := middleware.GetUserInfoFromContext(c)
	if user == nil {
		handler.Response(c, nil, nil)
		return
	}

	// 获取用户详细信息
	resp, err := ctrl.UserService().GetUserByUserId(c, user.UserId)
	resp.AccessToken = middleware.ExtractToken(c)

	// 设置访问 Cookie
	ctrl.SetVisitCookies(c, user.VisitToken, false)
	handler.Response(c, err, resp)
}
