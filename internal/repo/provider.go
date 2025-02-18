package repo

import "github.com/google/wire"

type Repo struct {
	ConfigRepo   ConfigRepo
	EmailRepo    EmailRepo
	SiteInfoRepo SiteInfoRepo
	UserRepo     UserRepo
}

func NewRepo(
	configRepo ConfigRepo,
	emailRepo EmailRepo,
	siteInfoRepo SiteInfoRepo,
	userRepo UserRepo,
) *Repo {
	return &Repo{
		ConfigRepo:   configRepo,
		EmailRepo:    emailRepo,
		SiteInfoRepo: siteInfoRepo,
		UserRepo:     userRepo,
	}
}

var ProviderSet = wire.NewSet(
	NewRepo,
	NewConfigRepo,
	NewEmailRepo,
	NewSiteInfoRepo,
	NewUserRepo,
)
