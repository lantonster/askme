package controller

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewController,
	NewUserController,
)

type Controller struct {
	UserController *UserController
}

func NewController(
	userController *UserController,
) *Controller {
	return &Controller{
		UserController: userController,
	}
}
