package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lantonster/askme/internal/model"
	"github.com/lantonster/askme/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestGetEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := service.NewConfigService(setupMock(ctrl))

	config := &model.Config{Email: &model.ConfigEmail{}}
	mockConfigRepo.EXPECT().FirstConfigByKey(context.Background(), model.ConfigKeyEmail).Return(config, nil)

	email, err := svc.GetEmail(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, config.Email, email)
}
