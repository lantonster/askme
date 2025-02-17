package model

import (
	"encoding/json"
	"time"

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

type SiteInfo struct {
	Id        int64
	Type      string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time

	Login *SiteInfoLogin `json:"login" gorm:"-"`
}

type SiteInfoLogin struct {
	AllowNewRegistrations   bool     `json:"allow_new_registrations"`
	AllowEmailRegistrations bool     `json:"allow_email_registrations"`
	AllowPasswordLogin      bool     `json:"allow_password_login"`
	LoginRequired           bool     `json:"login_required"`
	AllowEmailDomains       []string `json:"allow_email_domains"`
}

func (*SiteInfo) TableName() string {
	return "site_info"
}

func (si *SiteInfo) AfterFind(tx *gorm.DB) (err error) {
	switch si.Type {

	case string(SiteInfoTypeLogin):
		si.Login = &SiteInfoLogin{}
		json.Unmarshal([]byte(si.Content), si.Login)

	}

	return nil
}
