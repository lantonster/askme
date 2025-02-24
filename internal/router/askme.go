package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lantonster/askme/internal/controller"
)

type AskMeRouter struct {
	userController *controller.UserController
}

func NewAskMeRouter(
	userController *controller.UserController,
) *AskMeRouter {
	return &AskMeRouter{
		userController: userController,
	}
}

// 注册不需要鉴权的路由
func (r *AskMeRouter) RegisterNoAuth(router *gin.RouterGroup) {

	user := router.Group("/user")
	{

		// 邮箱确认
		user.POST("/email/verification", r.userController.VerifyEmail)

		// 当前用户信息
		user.GET("/info", r.userController.CurrentUserInfo)

		// 邮箱注册
		user.POST("/register/email", r.userController.RegisterUserByEmail)
	}

}
