package router

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lantonster/askme/docs"
	"github.com/lantonster/askme/internal/conf"
	"github.com/lantonster/askme/pkg/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type SwaggerRouter struct {
	swaggerConfig *conf.Swagger
}

func NewSwaggerRouter(config *conf.Config) *SwaggerRouter {
	return &SwaggerRouter{
		swaggerConfig: config.Swagger,
	}
}

func (r *SwaggerRouter) Register(engine *gin.RouterGroup) {
	if r.swaggerConfig.Show {
		docs.SwaggerInfo.Host = r.swaggerConfig.Host
		url := fmt.Sprintf("%s://%s/swagger/index.html", r.swaggerConfig.Protocal, r.swaggerConfig.Host)
		engine.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		log.WithContext(context.Background()).Infof("swagger 文档地址: %s", url)
	}
}
