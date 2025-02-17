package service

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewSiteInfoService,
	NewUploadsService,
	NewUserService,
)
