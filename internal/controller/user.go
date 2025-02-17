package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lantonster/askme/internal/middleware"
	"github.com/lantonster/askme/internal/schema"
	"github.com/lantonster/askme/internal/service"
	"github.com/lantonster/askme/pkg/errors"
	"github.com/lantonster/askme/pkg/handler"
	"github.com/lantonster/askme/pkg/i18n"
	"github.com/lantonster/askme/pkg/reason"
)

type UserController struct {
	*service.Service
}

func NewUserController(service *service.Service) *UserController {
	return &UserController{Service: service}
}

// RegisterUserByEmail godoc
//
//	@Summary		通过邮箱注册账号
//	@Description	通过邮箱注册账号
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			body	body		schema.RegisterUserByEmailReq						true	"body"
//	@Success		200		{object}	handler.ResponseBody{data=schema.RegisterUserByEmailRes}	"success"
//	@Router			/askme/api/v1/user/register/email [post]
func (uc *UserController) RegisterUserByEmail(c *gin.Context) {
	loginInfo, err := uc.SiteInfoService().GetSiteLogin(c)
	if err != nil {
		handler.Response(c, err, nil)
		return
	}
	if !loginInfo.AllowNewRegistrations || !loginInfo.AllowEmailRegistrations {
		handler.Response(c, errors.BadRequest(reason.UserRegistrationNotAllowed), nil)
		return
	}

	req := &schema.RegisterUserByEmailReq{}
	if handler.BindAndCheck(c, req) {
		return
	}
	if !loginInfo.IsEmailAllowed(req.Email) {
		handler.Response(c, errors.BadRequest(reason.EmailIllegalDomainError), nil)
		return
	}

	req.IP = c.ClientIP()
	if !middleware.GetUserIsAdminModerator(c) {
		// TODO
	}

	res, fieldErr, err := uc.UserService().RegisterUserByEmail(c, req)
	if len(fieldErr) > 0 {
		for _, field := range fieldErr {
			field.Error = i18n.Tr(handler.GetLang(c), field.Error)
		}
		handler.Response(c, err, fieldErr)
	} else {
		handler.Response(c, err, res)
	}
}
