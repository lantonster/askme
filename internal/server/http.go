package server

import (
	"context"
	"errors"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lantonster/askme/internal/conf"
	"github.com/lantonster/askme/internal/middleware"
	"github.com/lantonster/askme/internal/router"
	"github.com/lantonster/askme/pkg/log"
	"github.com/lantonster/askme/pkg/utils"
	"github.com/lantonster/askme/ui"
)

func NewHttpServer(
	config *conf.Config,
	router *router.Router,
	mid *middleware.Middleware,
) *Server {
	// 设置 gin 的运行模式
	gin.SetMode(utils.Ternary(config.Server.Http.Debug, gin.DebugMode, gin.ReleaseMode))

	engine := gin.New()

	// 注册 html 模板
	html, _ := fs.Sub(ui.Template, "template")
	htmlTemplate := template.Must(template.New("").Funcs(funcMap).ParseFS(html, "*")) // TODO funcmap
	engine.SetHTMLTemplate(htmlTemplate)

	// TODO middleware: langeuage, session, cors, logger, recovery, short id
	engine.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

	// 注册 UI 路由
	engine.Use(middleware.SetStaticCacheHeader) // 设置静态文件缓存
	router.Ui.Register(engine)

	// 注册 swagger 路由
	router.Swagger.Register(engine.Group("/swagger"))

	// 图片路由和登陆验证
	router.Uploads.Register(engine.Group("/uploads", mid.Avatar.AvatarThumb /* TODO vist auth */))

	// 注册 askme 路由
	askme := engine.Group("/askme/api/v1")
	{
		// 不需要鉴权
		router.AskMe.RegisterNoAuth(askme.Group("", mid.Auth.NoAuth))

		// 根据网站配置决定是否需要登陆
		askme.Group("", mid.Auth.NoAuth, mid.Auth.EjectUserBySiteInfo)

		// 需要登陆但不要求账号可用
		askme.Group("", mid.Auth.MustAuthWithoutAccountAvailable)

		// 需要登陆并且账号可用
		askme.Group("", mid.Auth.MustAuthAndAccountAvailable)

		// 管理端
	}

	return &Server{
		ShutdownTimeout: config.Server.Http.ShutdownTimeout,
		srv: &http.Server{
			Addr:    config.Server.Http.Addr,
			Handler: engine,
		},
	}
}

type Server struct {
	srv             *http.Server  // http 服务
	ShutdownTimeout time.Duration // 关闭超时时间
}

func (s *Server) Run(c context.Context) error {
	if s.srv == nil {
		return nil
	}

	quit := make(chan os.Signal, 1)
	errCh := make(chan error, 1)

	go func() {
		log.WithContext(c).Infof("http server start at %v", s.srv.Addr)
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.WithContext(c).Errorf("http server start failed: %v", err)
			errCh <- err
		}
	}()

	select {
	// 等待错误
	case err := <-errCh:
		s.Stop()
		return err

	// 等待退出信号
	case <-quit:
		return s.Stop()

	// 等待 c 关闭
	case <-c.Done():
		return s.Stop()
	}
}

func (s *Server) Stop() error {
	log.WithContext(context.Background()).Info("http server stop")
	c, cancel := context.WithTimeout(context.Background(), s.ShutdownTimeout)
	defer cancel()
	return s.srv.Shutdown(c)
}
