package service

import (
	"context"

	"github.com/lantonster/askme/internal/model"
	"github.com/lantonster/askme/internal/repo"
	"github.com/lantonster/askme/pkg/log"
)

type ConfigService interface {
	// GetConfigEmail 获取邮箱配置
	GetConfigEmail(c context.Context) (*model.ConfigEmail, error)
}

type ConfigServiceImpl struct {
	*repo.Repo
}

func NewConfigService(repo *repo.Repo) ConfigService {
	return &ConfigServiceImpl{Repo: repo}
}

func (s *ConfigServiceImpl) GetConfigEmail(c context.Context) (*model.ConfigEmail, error) {
	config, err := s.ConfigRepo.FirstConfigByKey(c, model.ConfigKeyEmail)
	if err != nil {
		log.WithContext(c).Errorf("获取邮箱配置失败: %v", err)
		return nil, err
	}
	return config.Email, nil
}
