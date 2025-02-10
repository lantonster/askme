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
	askmeRouter *router.AskMeRouter,
	swaggerRouter *router.SwaggerRouter,
	uiRouter *router.UiRouter,
	uploadsRouter *router.UploadsRouter,
	avatarMid *middleware.AvatarMiddleware,
) *Server {
	// 设置 gin 的运行模式
	gin.SetMode(utils.Ternary(config.Server.Http.Debug, gin.DebugMode, gin.ReleaseMode))

	r := gin.New()

	// 注册 html 模板
	html, _ := fs.Sub(ui.Template, "template")
	htmlTemplate := template.Must(template.New("").Funcs(funcMap).ParseFS(html, "*")) // TODO funcmap
	r.SetHTMLTemplate(htmlTemplate)

	// TODO middleware: langeuage, session, cors, logger, recovery, short id
	r.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

	// 注册 UI 路由
	r.Use(middleware.SetStaticCacheHeader) // 设置静态文件缓存
	uiRouter.Register(r)

	// 注册 swagger 路由
	swaggerRouter.Register(r.Group("/swagger"))

	// 图片路由和登陆验证
	uploadsRouter.Register(r.Group("/uploads", avatarMid.AvatarThumb /* TODO vist auth */))

	// 注册 askme 路由
	askme := r.Group("/askme/api/v1")
	{
		// 不需要鉴权的路由
		askmeRouter.RegisterNoAuth(askme)
	}

	return &Server{
		ShutdownTimeout: config.Server.Http.ShutdownTimeout,
		srv: &http.Server{
			Addr:    config.Server.Http.Addr,
			Handler: r,
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
