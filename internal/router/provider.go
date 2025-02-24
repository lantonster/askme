package router

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewRouter,
	NewAskMeRouter,
	NewSwaggerRouter,
	NewUiRouter,
	NewUploadsRouter,
)

type Router struct {
	AskMe   *AskMeRouter
	Swagger *SwaggerRouter
	Ui      *UiRouter
	Uploads *UploadsRouter
}

func NewRouter(
	askMe *AskMeRouter,
	swagger *SwaggerRouter,
	ui *UiRouter,
	uploads *UploadsRouter,
) *Router {
	return &Router{
		AskMe:   askMe,
		Swagger: swagger,
		Ui:      ui,
		Uploads: uploads,
	}
}
