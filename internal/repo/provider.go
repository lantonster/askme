package repo

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewRepo,
	NewActivityRepo,
	NewAuthRepo,
	NewConfigRepo,
	NewEmailRepo,
	NewSiteInfoRepo,
	NewUserRepo,
)

type Repo struct {
	ActivityRepo ActivityRepo
	AuthRepo     AuthRepo
	ConfigRepo   ConfigRepo
	EmailRepo    EmailRepo
	SiteInfoRepo SiteInfoRepo
	UserRepo     UserRepo
}

func NewRepo(
	activityRepo ActivityRepo,
	authRepo AuthRepo,
	configRepo ConfigRepo,
	emailRepo EmailRepo,
	siteInfoRepo SiteInfoRepo,
	userRepo UserRepo,
) *Repo {
	return &Repo{
		ActivityRepo: activityRepo,
		AuthRepo:     authRepo,
		ConfigRepo:   configRepo,
		EmailRepo:    emailRepo,
		SiteInfoRepo: siteInfoRepo,
		UserRepo:     userRepo,
	}
}
