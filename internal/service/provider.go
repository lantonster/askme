package service

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewService,
	NewSiteInfoService,
	NewUploadsService,
	NewUserService,
)

var (
	siteInfoService SiteInfoService
	uploadsService  UploadsService
	userService     UserService
)

type Service struct{}

func NewService(
	siteInfo SiteInfoService,
	uploads UploadsService,
	user UserService,
) *Service {
	siteInfoService = siteInfo
	uploadsService = uploads
	userService = user

	return &Service{}
}

func (s *Service) SiteInfoService() SiteInfoService {
	return siteInfoService
}

func (s *Service) UploadsService() UploadsService {
	return uploadsService
}

func (s *Service) UserService() UserService {
	return userService
}
