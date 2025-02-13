package service

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/lantonster/askme/internal/conf"
	"github.com/lantonster/askme/pkg/dir"
	"github.com/lantonster/askme/pkg/errors"
	"github.com/lantonster/askme/pkg/reason"
)

var (
	// 支持的缩略图文件后缀
	supportedThumbFileExtMapping = map[string]imaging.Format{
		".jpg":  imaging.JPEG,
		".jpeg": imaging.JPEG,
		".png":  imaging.PNG,
		".gif":  imaging.GIF,
	}
)

type UploadsService interface {
	AvatarThumbFile(c *gin.Context, fileName string, size int) (url string, err error)
}

type uploadsService struct {
	uploadsConfig *conf.Uploads
}

func NewUploadsService(config *conf.Config) UploadsService {
	return &uploadsService{uploadsConfig: config.Uploads}
}

// AvatarThumbFile 生成头像缩略图文件并返回其路径
func (s *uploadsService) AvatarThumbFile(c *gin.Context, fileName string, size int) (url string, err error) {
	suffix := path.Ext(fileName)                                  // 获取文件后缀
	filePath := path.Join(s.uploadsConfig.AvatarPath(), fileName) // 获取文件路径

	// 检查目标文件类型是否支持压缩，如果不支持直接返回保存路径
	if _, ok := supportedThumbFileExtMapping[suffix]; !ok {
		return filePath, nil
	}

	// 限制最大压缩尺寸
	if size > 1024 {
		size = 1024
	}

	thumbFileName := fmt.Sprintf("%d_%d@%s", size, size, fileName)
	thumbFilePath := path.Join(s.uploadsConfig.AvatarThumbSubPath(), thumbFileName)

	// 检查压缩文件是否存在，如果存在返回对应的保存路径
	if _, err := os.ReadFile(thumbFilePath); err == nil {
		return thumbFilePath, nil
	}

	// 读取原头像文件
	avatarFile, err := os.ReadFile(filePath)
	if err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithMsg("读取头像文件 %s 时出错", filePath).WithError(err).WithStack()
	}

	// 解码头像文件
	reader := bytes.NewReader(avatarFile)
	img, err := imaging.Decode(reader)
	if err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithMsg("解码头像文件 %s 时出错", filePath).WithError(err).WithStack()
	}

	// 压缩并编码头像文件
	var buf bytes.Buffer
	newImg := imaging.Fill(img, size, size, imaging.Center, imaging.Linear)
	if err := imaging.Encode(&buf, newImg, supportedThumbFileExtMapping[suffix]); err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithMsg("压缩后编码头像文件 %s 时出错", filePath).WithError(err).WithStack()
	}

	// 创建头像缩略图目录
	if err := dir.CreateDirIfNotExist(s.uploadsConfig.AvatarThumbSubPath()); err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithMsg("创建头像缩略图目录 %s 时出错", s.uploadsConfig.AvatarThumbSubPath()).WithError(err).WithStack()
	}

	// 创建头像缩略图文件
	out, err := os.Create(thumbFilePath)
	if err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithMsg("创建头像缩略图文件 %s 时出错", thumbFilePath).WithError(err).WithStack()
	}
	defer out.Close()

	// 保存头像缩略图文件
	thumbReader := bytes.NewReader(buf.Bytes())
	if _, err := io.Copy(out, thumbReader); err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithMsg("保存头像缩略图文件 %s 时出错", thumbFilePath).WithError(err).WithStack()
	}

	return thumbFilePath, nil
}
