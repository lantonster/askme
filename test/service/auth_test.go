package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lantonster/askme/internal/model"
	"github.com/lantonster/askme/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestSetUserCacheInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := service.NewAuthService(setupMock(ctrl))

	mockAuthRepo.EXPECT().SetUserCache(context.Background(), gomock.Any(), gomock.Any()).Return(nil)

	accessToken, visitToken, err := svc.SetUserCacheInfo(context.Background(), &model.User{})
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, visitToken)
}
