package router

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewAskMeRouter,
	NewSwaggerRouter,
	NewUiRouter,
	NewUploadsRouter,
)
