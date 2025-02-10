package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lantonster/askme/internal/schema"
	"github.com/lantonster/askme/pkg/handler"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
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
	// TODO check site

	req := &schema.RegisterUserByEmailReq{}
	if handler.BindAndCheck(c, req) {
		return
	}
}
