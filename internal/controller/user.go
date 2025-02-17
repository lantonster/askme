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
	siteInfoService service.SiteInfoService
}

func NewUserController(
	siteInfoService service.SiteInfoService,
) *UserController {
	return &UserController{
		siteInfoService: siteInfoService,
	}
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
	siteInfo, err := uc.siteInfoService.GetSiteLogin(c)
	if err != nil {
		handler.Response(c, err, nil)
		return
	}
	if !siteInfo.AllowNewRegistrations || !siteInfo.AllowEmailRegistrations {
		handler.Response(c, errors.BadRequest(reason.UserRegistrationNotAllowed), nil)
		return
	}

	req := &schema.RegisterUserByEmailReq{}
	if handler.BindAndCheck(c, req) {
		return
	}
}
