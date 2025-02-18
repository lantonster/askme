package model

import (
	"encoding/json"

	"github.com/lantonster/askme/pkg/gravatar"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

const (
	// UserStatusAvailable 用户可用
	UserStatusAvailable = "available"
	// UserStatusSuspended 用户被暂停(禁言)
	UserStatusSuspended = "suspended"
	// UserStatusDeleted 用户已删除
	UserStatusDeleted = "deleted"
	// UserStatusInactive 用户未激活
	UserStatusInactive = "inactive"
)

const (
	// EmailStatusAvailable 邮件可用
	EmailStatusAvailable = "available"
	// EmailStatusToBeVerified 邮件待验证
	EmailStatusToBeVerified = "to_be_verified"
)

// User [MySQL] 用户
type User struct {
	Id            int64
	CreatedAt     int64
	UpdatedAt     int64
	DeletedAt     soft_delete.DeletedAt
	SuspendedAt   int64 // 被禁言的时间
	LastLoginDate int64 // 上次登录时间

	Username    string // 用户名
	Password    string // 密码
	Email       string // 邮箱
	DisplayName string // 昵称
	IP          string // IP 地址
	Status      string // 用户状态
	MailStatus  string // 邮箱状态
	Avatar      string // 头像
}

func (User) TableName() string {
	return "user"
}

func (u *User) AfterFind(tx *gorm.DB) error {
	// 如果用户的邮件状态为待验证，将用户状态设置为未激活
	if u.MailStatus == EmailStatusToBeVerified {
		u.Status = UserStatusInactive
	}

	avatar := &AvatarInfo{}
	json.Unmarshal([]byte(u.Avatar), avatar)

	// 如果用户被删除，将用户头像设置为空
	if u.Status == UserStatusDeleted {
		u.Avatar = ""
	}
	switch avatar.Type {
	case AvatarTypeGravatar:
		u.Avatar = gravatar.GetAvatarURL(DefaultGravatarBaseURL, u.Email)
	case AvatarTypeCustom:
		u.Avatar = avatar.Custom
	default:
		u.Avatar = ""
	}

	return nil
}

const (
	// 默认头像配置
	DefaultGravatarBaseURL = "https://www.gravatar.com/avatar/"
	// 默认头像
	DefaultAvatar = "system"
	// 默认头像类型
	AvatarTypeDefault = "default"
	// Gravatar 头像
	AvatarTypeGravatar = "gravatar"
	// 自定义头像
	AvatarTypeCustom = "custom"
)

type AvatarInfo struct {
	Type     string `json:"type"`
	Gravatar string `json:"gravatar"`
	Custom   string `json:"custom"`
}
