package model

type UserCache struct {
	Id         int64  `json:"id"`
	UserStatus string `json:"user_status"`
	MailStatus string `json:"mail_status"`
}
