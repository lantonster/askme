package repo

import "github.com/google/wire"

type Repo struct {
	AuthRepo     AuthRepo
	ConfigRepo   ConfigRepo
	EmailRepo    EmailRepo
	SiteInfoRepo SiteInfoRepo
	UserRepo     UserRepo
}

func NewRepo(
	authRepo AuthRepo,
	configRepo ConfigRepo,
	emailRepo EmailRepo,
	siteInfoRepo SiteInfoRepo,
	userRepo UserRepo,
) *Repo {
	return &Repo{
		AuthRepo:     authRepo,
		ConfigRepo:   configRepo,
		EmailRepo:    emailRepo,
		SiteInfoRepo: siteInfoRepo,
		UserRepo:     userRepo,
	}
}

var ProviderSet = wire.NewSet(
	NewRepo,
	NewAuthRepo,
	NewConfigRepo,
	NewEmailRepo,
	NewSiteInfoRepo,
	NewUserRepo,
)
