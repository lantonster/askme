package router

import (
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lantonster/askme/internal/conf"
)

type UploadsRouter struct {
	uploadsConfig *conf.Uploads
}

func NewUploadsRouter(config *conf.Config) *UploadsRouter {
	return &UploadsRouter{
		uploadsConfig: config.Uploads,
	}
}

func (r *UploadsRouter) Register(engine *gin.RouterGroup) {
	engine.Static("/avatar", r.uploadsConfig.AvatarPath())
	engine.Static("/avatar_thumb", r.uploadsConfig.AvatarThumbSubPath())
	engine.Static("/post", r.uploadsConfig.PostPath())
	engine.Static("/branding", r.uploadsConfig.BrandingPath())
	engine.GET("/files/post/*filepath", func(c *gin.Context) {
		// 获取请求中的文件路径参数，例如 hash/123.pdf
		filePath := c.Param("filepath")
		// 获取原始文件名，例如 123.pdf
		originalFilename := filepath.Base(filePath)
		// 生成实际文件名，例如 hash.pdf
		realFilename := strings.TrimSuffix(filePath, "/"+originalFilename) + filepath.Ext(originalFilename)
		// 生成文件在本地的路径，例如 /uploads/files/post/hash.pdf
		fileLocalPath := filepath.Join(r.uploadsConfig.FilePostPath(), realFilename)
		// 将文件作为附件返回，供下载
		c.FileAttachment(fileLocalPath, originalFilename)
	})
}
