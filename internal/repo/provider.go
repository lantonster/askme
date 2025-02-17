package repo

import "github.com/google/wire"

type Repo struct {
	SiteInfo SiteInfoRepo
	User     UserRepo
}

func NewRepo(
	siteInfo SiteInfoRepo,
	user UserRepo,
) *Repo {
	return &Repo{
		SiteInfo: siteInfo,
		User:     user,
	}
}

var ProviderSet = wire.NewSet(
	NewRepo,
	NewSiteInfoRepo,
	NewUserRepo,
)
