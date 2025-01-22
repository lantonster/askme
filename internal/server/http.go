package server

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lantonster/askme/internal/conf"
	"github.com/lantonster/askme/internal/router"
	"github.com/lantonster/askme/pkg/utils"
	"github.com/lantonster/askme/ui"
)

func NewHttpServer(
	config *conf.Config,
	uiRouter *router.UiRouter,
) *Server {
	// 设置 gin 的运行模式
	gin.SetMode(utils.Ternary(config.Server.Http.Debug, gin.DebugMode, gin.ReleaseMode))

	r := gin.New()
	r.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

	// 注册 html 模板
	html, _ := fs.Sub(ui.Template, "template")
	htmlTemplate := template.Must(template.New("").Funcs(funcMap).ParseFS(html, "*"))
	r.SetHTMLTemplate(htmlTemplate)

	// 注册 UI 路由
	uiRouter.Register(r, config)

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
		fmt.Println("http server start at", s.srv.Addr)
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("http server start failed:", err)
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
	fmt.Println("http server stop")
	c, cancel := context.WithTimeout(context.Background(), s.ShutdownTimeout)
	defer cancel()
	return s.srv.Shutdown(c)
}
