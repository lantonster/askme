package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lantonster/askme/internal/schema"
	"github.com/lantonster/askme/internal/service"
	"github.com/lantonster/askme/pkg/errors"
	"github.com/lantonster/askme/pkg/handler"
	"github.com/lantonster/askme/pkg/reason"
)

type UserController struct {
	*service.Service
}

func NewUserController(service *service.Service) *UserController {
	return &UserController{Service: service}
}

// VerifyEmail godoc
//
//	@Summary		邮箱验证
//	@Description	邮箱验证
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			body	body		schema.VerificationEmailReq								true	"body"
//	@Success		200		{object}	handler.ResponseBody{data=schema.VerificationEmailRes}	"success"
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

	// TODO action del
	handler.Response(c, nil, resp)
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
//	@Router			/askme/api/v1/user/register/email [post]
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
