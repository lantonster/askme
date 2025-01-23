package middleware

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lantonster/askme/internal/conf"
	"github.com/lantonster/askme/internal/service/uploads"
	"github.com/spf13/cast"
)

type AvatarMiddleware struct {
	uploadsConfig  *conf.Uploads
	uploadsService uploads.UploadsService
}

func NewAvatarMiddleware(config *conf.Config, uploadsService uploads.UploadsService) *AvatarMiddleware {
	return &AvatarMiddleware{
		uploadsConfig:  config.Uploads,
		uploadsService: uploadsService,
	}
}

// AvatarThumb 头像缩略图
func (a *AvatarMiddleware) AvatarThumb(c *gin.Context) {
	uri := c.Request.RequestURI
	uriWithoutQuery, err := url.Parse(uri)
	if err != nil {
		fmt.Println("解析 URL 失败", err)
		c.Next()
		return
	}
	ext := strings.TrimPrefix(path.Ext(uriWithoutQuery.Path), ".")
	c.Header("content-type", fmt.Sprintf("image/%s", ext))

	if strings.HasPrefix(uri, "/uploads/avatar/") {
		size := cast.ToInt(c.Query("s"))
		filename := filepath.Base(uriWithoutQuery.Path)
		filepath := path.Join(a.uploadsConfig.AvatarPath(), filename)

		var err error
		if size != 0 {
			// 获取头像缩略图地址
			filepath, err = a.uploadsService.AvatarThumbFile(c, filename, size)
			if err != nil {
				fmt.Println("获取头像缩略图地址失败", err)
				c.Abort()
				return
			}
		}

		// 读取头像缩略图文件
		avatarFile, err := os.ReadFile(filepath)
		if err != nil {
			fmt.Println("读取头像缩略图文件失败", err)
			c.Abort()
			return
		}

		// 写入头像缩略图文件
		if _, err := c.Writer.Write(avatarFile); err != nil {
			fmt.Println("写入头像缩略图文件失败", err)
		}
		c.Abort()
		return
	}

	c.Next()
}
