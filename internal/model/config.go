package model

import (
	"encoding/json"

	"gorm.io/gorm"
)

const (
	ConfigKeyEmail = "email"
)

// Config [MySQL / Redis] 配置
type Config struct {
	Id    int64
	Key   string
	Value string

	Email *ConfigEmail `gorm:"-"`
}

func (Config) TableName() string {
	return "config"
}

func (c *Config) AfterFind(tx *gorm.DB) error {
	switch c.Key {

	case ConfigKeyEmail:
		c.Email = &ConfigEmail{}
		json.Unmarshal([]byte(c.Value), c.Email)

	}
	return nil
}

type ConfigEmail struct {
	FromName           string `json:"from_name"`
	FromEmail          string `json:"from_email"`
	SMTPHost           string `json:"smtp_host"`
	SMTPPort           int    `json:"smtp_port"`
	SMTPPassword       string `json:"smtp_password"`
	SMTPUsername       string `json:"smtp_username"`
	SMTPAuthentication bool   `json:"smtp_authentication"`
	Encryption         string `json:"encryption"`
	RegisterTitle      string `json:"register_title"`
	RegisterBody       string `json:"register_body"`
	PassResetTitle     string `json:"pass_reset_title"`
	PassResetBody      string `json:"pass_reset_body"`
	ChangeTitle        string `json:"change_title"`
	ChangeBody         string `json:"change_body"`
	TestTitle          string `json:"test_title"`
	TestBody           string `json:"test_body"`
	NewAnswerTitle     string `json:"new_answer_title"`
	NewAnswerBody      string `json:"new_answer_body"`
	NewCommentTitle    string `json:"new_comment_title"`
	NewCommentBody     string `json:"new_comment_body"`
}

func (c *ConfigEmail) IsSSL() bool {
	return c.Encryption == "SSL"
}

func (c *ConfigEmail) IsTLS() bool {
	return c.Encryption == "TLS"
}
