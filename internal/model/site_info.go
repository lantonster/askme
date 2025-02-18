package model

import (
	"encoding/json"
	"strings"

	"gorm.io/gorm"
)

type SiteInfoType string

const (
	SiteInfoTypeGeneral       SiteInfoType = "general"
	SiteInfoTypeInterface     SiteInfoType = "interface"
	SiteInfoTypeBranding      SiteInfoType = "branding"
	SiteInfoTypeWrite         SiteInfoType = "write"
	SiteInfoTypeLegal         SiteInfoType = "legal"
	SiteInfoTypeSeo           SiteInfoType = "seo"
	SiteInfoTypeLogin         SiteInfoType = "login"
	SiteInfoTypeCustomCssHTML SiteInfoType = "css-html"
	SiteInfoTypeTheme         SiteInfoType = "theme"
	SiteInfoTypePrivileges    SiteInfoType = "privileges"
	SiteInfoTypeUsers         SiteInfoType = "users"
)

// SiteInfo [MySQL / Redis] 网站信息
type SiteInfo struct {
	Id        int64
	Type      string
	Content   string
	CreatedAt int64
	UpdatedAt int64

	Genral *SiteInfoGeneral `json:"general" gorm:"-"`
	Login  *SiteInfoLogin   `json:"login" gorm:"-"`
	Users  *SiteInfoUsers   `json:"users" gorm:"-"`
}

func (*SiteInfo) TableName() string {
	return "site_info"
}

func (si *SiteInfo) AfterFind(tx *gorm.DB) (err error) {
	switch si.Type {

	case string(SiteInfoTypeGeneral):
		si.Genral = &SiteInfoGeneral{}
		json.Unmarshal([]byte(si.Content), si.Genral)

	case string(SiteInfoTypeLogin):
		si.Login = &SiteInfoLogin{}
		json.Unmarshal([]byte(si.Content), si.Login)

	case string(SiteInfoTypeUsers):
		si.Users = &SiteInfoUsers{}
		json.Unmarshal([]byte(si.Content), si.Users)

	}

	return nil
}

type SiteInfoGeneral struct {
	CheckUpdate      bool   `json:"check_update"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	Description      string `json:"description"`
	SiteUrl          string `json:"site_url"`
	ContactEmail     string `json:"contact_email"`
}

type SiteInfoLogin struct {
	AllowNewRegistrations   bool     `json:"allow_new_registrations"`
	AllowEmailRegistrations bool     `json:"allow_email_registrations"`
	AllowPasswordLogin      bool     `json:"allow_password_login"`
	LoginRequired           bool     `json:"login_required"`
	AllowEmailDomains       []string `json:"allow_email_domains"`
}

// IsEmailAllowed 函数用于检查给定的电子邮件是否在当前登录配置的允许域名列表中。
//
// 参数:
//   - email: 要检查的电子邮件地址
//
// 返回: 如果电子邮件在允许列表中则返回 true，否则返回 false
func (login *SiteInfoLogin) IsEmailAllowed(email string) bool {
	if len(login.AllowEmailDomains) == 0 {
		return true
	}

	for _, domain := range login.AllowEmailDomains {
		if strings.HasSuffix(email, domain) {
			return true
		}
	}
	return false
}

type SiteInfoUsers struct {
	AllowUpdateAvatar      bool   `json:"allow_update_avatar"`
	AllowUpdateBio         bool   `json:"allow_update_bio"`
	AllowUpdateDisplayName bool   `json:"allow_update_display_name"`
	AllowUpdateLocation    bool   `json:"allow_update_location"`
	AllowUpdateUsername    bool   `json:"allow_update_username"`
	AllowUpdateWebsite     bool   `json:"allow_update_website"`
	DefaultAvatar          string `json:"default_avatar"`
	GravatarBaseURL        string `json:"gravatar_base_url"`
}
