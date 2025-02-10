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

	// TODO

	user := router.Group("/user")
	{
		user.POST("/register/email", r.userController.RegisterUserByEmail)
	}

}
