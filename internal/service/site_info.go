package service

import (
	"context"

	"github.com/lantonster/askme/internal/model"
	"github.com/lantonster/askme/internal/repo"
	"github.com/lantonster/askme/internal/schema"
	"github.com/lantonster/askme/pkg/log"
)

type SiteInfoService interface {
	// GetSiteGeneral 获取站点常规信息
	GetSiteGeneral(c context.Context) (*schema.GetSiteGeneralRes, error)

	// GetSiteLogin 获取站点登录信息
	GetSiteLogin(c context.Context) (*schema.GetSiteLoginRes, error)

	// GetSiteUsers 获取站点用户信息
	GetSiteUsers(c context.Context) (*schema.GetSiteUsersRes, error)
}

type SiteInfoServiceImpl struct {
	*repo.Repo
}

func NewSiteInfoService(repo *repo.Repo) SiteInfoService {
	return &SiteInfoServiceImpl{Repo: repo}
}

func (s *SiteInfoServiceImpl) GetSiteGeneral(c context.Context) (*schema.GetSiteGeneralRes, error) {
	siteInfo, err := s.SiteInfoRepo.FirstSiteInfoByType(c, model.SiteInfoTypeGeneral)
	if err != nil {
		log.WithContext(c).Errorf("获取站点常规信息失败: %v", err)
		return nil, err
	}
	return siteInfo.Genral, nil
}

func (s *SiteInfoServiceImpl) GetSiteLogin(c context.Context) (*schema.GetSiteLoginRes, error) {
	siteInfo, err := s.SiteInfoRepo.FirstSiteInfoByType(c, model.SiteInfoTypeLogin)
	if err != nil {
		log.WithContext(c).Errorf("获取站点登录信息失败: %v", err)
		return nil, err
	}
	return siteInfo.Login, nil
}

func (s *SiteInfoServiceImpl) GetSiteUsers(c context.Context) (*schema.GetSiteUsersRes, error) {
	siteInfo, err := s.SiteInfoRepo.FirstSiteInfoByType(c, model.SiteInfoTypeUsers)
	if err != nil {
		log.WithContext(c).Errorf("获取站点用户信息失败: %v", err)
		return nil, err
	}
	return siteInfo.Users, nil
}
