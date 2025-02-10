//go:build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/lantonster/askme/internal/conf"
	"github.com/lantonster/askme/internal/controller"
	"github.com/lantonster/askme/internal/middleware"
	"github.com/lantonster/askme/internal/router"
	"github.com/lantonster/askme/internal/server"
	"github.com/lantonster/askme/internal/service"
)

func Init() *server.Server {
	panic(wire.Build(
		conf.ProviderSet,
		server.ProviderSet,
		router.ProviderSet,
		service.ProviderSet,
		middleware.ProviderSet,
		controller.ProviderSet,
	))
}
