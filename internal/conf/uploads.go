package conf

import (
	"path"
)

type Uploads struct {
	Path string `yaml:"path" mapstructure:"path"` // 上传路径
}

// AvatarPath 头像上传路径
func (u *Uploads) AvatarPath() string {
	return path.Join(u.Path, "avatar")
}

// AvatarThumbPath 头像缩略图上传路径
func (u *Uploads) AvatarThumbSubPath() string {
	return path.Join(u.Path, "avatar_thumb")
}

func (u *Uploads) PostPath() string {
	return path.Join(u.Path, "post")
}

func (u *Uploads) BrandingPath() string {
	return path.Join(u.Path, "branding")
}

func (u *Uploads) FilePostPath() string {
	return path.Join(u.Path, "files/post")
}
