package router

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lantonster/askme/internal/conf"
	"github.com/lantonster/askme/pkg/log"
	"github.com/lantonster/askme/ui"
)

type UiRouter struct {
	uiConfig *conf.UI
}

func NewUiRouter(config *conf.Config) *UiRouter {
	return &UiRouter{
		uiConfig: config.UI,
	}
}

type _resource struct {
	fs embed.FS
}

// Open 实现 http.FileSystem 接口，用于获取静态资源。
//
// 这里的 name 是请求路径中的静态资源路径，即请求 {baseUrl}/static/xxx 时，name 为 xxx。
// 这里的 name 是相对于 ui/build/static 路径的，所以需要拼接上 static
func (r *_resource) Open(name string) (fs.File, error) {
	name = fmt.Sprintf("build/static/%s", name)
	return r.fs.Open(name)
}

func (r *UiRouter) Register(engine *gin.Engine) {
	baseUrl := r.uiConfig.BaseUrl

	// 注册静态资源路径为 {baseUrl}/static，即请求 {baseUrl}/static/xxx 时，会从 ui/build/static/xxx 路径下获取资源
	engine.StaticFS(baseUrl+"/static", http.FS(&_resource{fs: ui.Build}))

	engine.NoRoute(func(c *gin.Context) {
		urlPath := c.Request.URL.Path

		// 如果设置的 baseUrl 不为空，则需要去对请求路径去除 baseUrl 前缀
		if len(baseUrl) > 0 {
			urlPath = strings.TrimPrefix(urlPath, baseUrl)
		}

		var filePath string
		switch urlPath {
		case "/favicon.ico":
			// branding, err := a.siteInfoService.GetSiteBranding(c)
			// if err != nil {
			// 	log.Error(err)
			// }
			// if branding.Favicon != "" {
			// 	c.String(http.StatusOK, htmltext.GetPicByUrl(branding.Favicon))
			// 	return
			// } else if branding.SquareIcon != "" {
			// 	c.String(http.StatusOK, htmltext.GetPicByUrl(branding.SquareIcon))
			// 	return
			// } else {
			// 	c.Header("content-type", "image/vnd.microsoft.icon")
			// 	filePath = UIRootFilePath + urlPath
			// }
			filePath = "build/favicon.ico"
		case "/manifest.json":
			// a.siteInfoController.GetManifestJson(c)
			// return
			filePath = "build/manifest.json"
		case "/install":
			// 如果是通过命令行运行的，则无法访问安装页面
			c.Redirect(http.StatusFound, "/")
			return
		default:
			filePath = "build/index.html"
			c.Header("content-type", "text/html;charset=utf-8")
			c.Header("X-Frame-Options", "DENY")
		}

		file, err := ui.Build.ReadFile(filePath)
		if err != nil {
			log.WithContext(c).Errorf("读取文件 %s 失败: %v", filePath, err)
			c.Status(http.StatusNotFound)
			return
		}

		// cdnPrefix := ""
		// _ = plugin.CallCDN(func(fn plugin.CDN) error {
		// 	cdnPrefix = fn.GetStaticPrefix()
		// 	return nil
		// })
		// if cdnPrefix != "" {
		// 	if cdnPrefix[len(cdnPrefix)-1:] == "/" {
		// 		cdnPrefix = strings.TrimSuffix(cdnPrefix, "/")
		// 	}
		// 	c.String(http.StatusOK, strings.ReplaceAll(string(file), "/static", cdnPrefix+"/static"))
		// 	return
		// }

		// This part is to solve the problem of returning 404 when the access path does not exist.
		// However, there is no way to check whether the current route exists in the frontend.
		// We can only hand over the route to the frontend for processing.
		// And the plugin, frontend routes can now be dynamically registered,
		// so there's no good way to get all frontend routes
		//if filePath == UIIndexFilePath {
		//	c.String(http.StatusNotFound, string(file))
		//	return
		//}

		c.String(http.StatusOK, string(file))
	})

}
