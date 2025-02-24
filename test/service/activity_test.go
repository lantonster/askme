package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lantonster/askme/internal/model"
	"github.com/lantonster/askme/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestActivateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := service.NewActivityService(setupMock(ctrl))

	user := &model.User{Id: 1, Rank: 100}
	config := &model.Config{Id: 1, Key: model.ConfigKeyUserActivated, Value: "1", UserActivated: 1}

	t.Run("查询用户激活配置错误", func(t *testing.T) {
		targetErr := fmt.Errorf("first config error")
		mockConfigRepo.EXPECT().FirstConfigByKey(context.Background(), model.ConfigKeyUserActivated).Return(nil, targetErr)

		err := svc.ActivateUser(context.Background(), user)
		assert.Equal(t, targetErr, err)
	})

	t.Run("激活用户活动已存在", func(t *testing.T) {
		mockConfigRepo.EXPECT().FirstConfigByKey(context.Background(), model.ConfigKeyUserActivated).Return(config, nil)
		mockActivityRepo.EXPECT().GetActivityByType(context.Background(), model.ActivityTypeUserActivated).Return(nil, true, nil)

		err := svc.ActivateUser(context.Background(), user)
		assert.NoError(t, err)
	})

	t.Run("保存激活用户活动成功", func(t *testing.T) {
		mockConfigRepo.EXPECT().FirstConfigByKey(context.Background(), model.ConfigKeyUserActivated).Return(config, nil)
		mockActivityRepo.EXPECT().GetActivityByType(context.Background(), model.ActivityTypeUserActivated).Return(nil, false, nil)
		mockUserRepo.EXPECT().IncrRank(context.Background(), user.Id, user.Rank, config.UserActivated).Return(nil)
		mockActivityRepo.EXPECT().CreateActivity(context.Background(), gomock.Any()).Return(nil)

		err := svc.ActivateUser(context.Background(), user)
		assert.NoError(t, err)
	})
}
