package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

const (
	UserStatusAvailable = "available"
	UserStatusSuspended = "suspended"
	UserStatusDeleted   = "deleted"
)

const (
	EmailStatusAvailable    = "available"
	EmailStatusToBeVerified = "to_be_verified"
)

// User [MySQL] 用户
type User struct {
	Id            int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
	SuspendedAt   sql.NullTime // 被禁言的时间
	LastLoginDate sql.NullTime // 上次登录时间

	Username    string // 用户名
	Password    string // 密码
	Email       string // 邮箱
	DisplayName string // 昵称
	IP          string // IP 地址
	Status      string // 用户状态
	MailStatus  string // 邮箱状态
}

func (User) TableName() string {
	return "user"
}
