package service

import (
	"context"

	"github.com/lantonster/askme/internal/model"
	"github.com/lantonster/askme/internal/repo"
	"github.com/lantonster/askme/internal/schema"
	"github.com/lantonster/askme/pkg/log"
)

type SiteInfoService interface {
	// GetSiteLogin 获取站点登录信息
	GetSiteLogin(c context.Context) (*schema.GetSiteLoginRes, error)
}

type siteInfoService struct {
	siteInfoRepo repo.SiteInfoRepo
}

func NewSiteInfoService(siteInfoRepo repo.SiteInfoRepo) SiteInfoService {
	return &siteInfoService{
		siteInfoRepo: siteInfoRepo,
	}
}

func (s *siteInfoService) GetSiteLogin(c context.Context) (*schema.GetSiteLoginRes, error) {
	siteInfo, err := s.siteInfoRepo.FirstSiteInfoByType(c, model.SiteInfoTypeLogin)
	if err != nil {
		log.WithContext(c).Errorf("获取站点登录信息失败: %v", err)
		return nil, err
	}

	return siteInfo.Login, nil
}
