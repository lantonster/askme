package service

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewService,
	NewAuthService,
	NewConfigService,
	NewEmailService,
	NewSiteInfoService,
	NewUploadsService,
	NewUserService,
)

var (
	authService     AuthService
	configService   ConfigService
	emailService    EmailService
	siteInfoService SiteInfoService
	uploadsService  UploadsService
	userService     UserService
)

type Service struct{}

func NewService(
	auth AuthService,
	config ConfigService,
	email EmailService,
	siteInfo SiteInfoService,
	uploads UploadsService,
	user UserService,
) *Service {
	authService = auth
	configService = config
	emailService = email
	siteInfoService = siteInfo
	uploadsService = uploads
	userService = user

	return &Service{}
}

func (s *Service) Auth() AuthService {
	return authService
}

func (s *Service) ConfigService() ConfigService {
	return configService
}

func (s *Service) EmailService() EmailService {
	return emailService
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
