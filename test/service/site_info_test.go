package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lantonster/askme/internal/model"
	"github.com/lantonster/askme/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestGetSiteGeneral(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := service.NewSiteInfoService(setupMock(ctrl))

	config := &model.SiteInfo{Genral: &model.SiteInfoGeneral{}}
	mockSiteInfoRepo.EXPECT().FirstSiteInfoByType(context.Background(), model.SiteInfoTypeGeneral).Return(config, nil)

	general, err := svc.GetSiteGeneral(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, config.Genral, general)
}

func TestGetSiteLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := service.NewSiteInfoService(setupMock(ctrl))

	config := &model.SiteInfo{Login: &model.SiteInfoLogin{}}
	mockSiteInfoRepo.EXPECT().FirstSiteInfoByType(context.Background(), model.SiteInfoTypeLogin).Return(config, nil)

	login, err := svc.GetSiteLogin(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, config.Login, login)
}

func TestGetSiteUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := service.NewSiteInfoService(setupMock(ctrl))

	config := &model.SiteInfo{Users: &model.SiteInfoUsers{}}
	mockSiteInfoRepo.EXPECT().FirstSiteInfoByType(context.Background(), model.SiteInfoTypeUsers).Return(config, nil)

	users, err := svc.GetSiteUsers(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, config.Users, users)
}
