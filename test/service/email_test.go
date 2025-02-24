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

func TestSend(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := service.NewEmailService(setupMock(ctrl))

	t.Run("获取邮件配置失败", func(t *testing.T) {
		targetErr := fmt.Errorf("get email err")
		mockConfigService.EXPECT().GetEmail(context.Background()).Return(nil, targetErr)

		err := svc.Send(context.Background(), "342310798@qq.com", "123", "123")
		assert.Equal(t, targetErr, err)
	})

	t.Run("发送邮件", func(t *testing.T) {
		email := &model.ConfigEmail{
			FromName:     "lllllan",
			FromEmail:    "342310798@qq.com",
			SMTPHost:     "smtp.qq.com",
			SMTPPort:     465,
			SMTPPassword: "lgnjnekdmnaacagd",
			SMTPUsername: "342310798@qq.com",
			Encryption:   "ssl",
		}

		mockConfigService.EXPECT().GetEmail(context.Background()).Return(email, nil)

		err := svc.Send(context.Background(), "342310798@qq.com", "title", "body")
		assert.NoError(t, err)
	})
}

func TestSendRegisterVerificationEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := service.NewEmailService(setupMock(ctrl))

	t.Run("获取站点通用信息失败", func(t *testing.T) {
		var targetErr = fmt.Errorf("get site general err")
		mockSiteInfoService.EXPECT().GetSiteGeneral(context.Background()).Return(nil, targetErr)

		err := svc.SendRegisterVerificationEmail(context.Background(), 1, "test@example.com")
		fmt.Printf("err: %v\n", err)
	})

	t.Run("发送验证码", func(t *testing.T) {
		general := &model.SiteInfoGeneral{}
		mockSiteInfoService.EXPECT().GetSiteGeneral(context.Background()).Return(general, nil)
		mockEmailRepo.EXPECT().StoreVerificationEmail(context.Background(), int64(1), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		mockEmailService.EXPECT().Send(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		err := svc.SendRegisterVerificationEmail(context.Background(), 1, "test@example.com")
		assert.NoError(t, err)
	})
}

func TestVerifyUrlExpired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := service.NewEmailService(setupMock(ctrl))

	t.Run("验证码过期", func(t *testing.T) {
		code := "code"
		mockEmailRepo.EXPECT().VerifyCode(context.Background(), code).Return(nil, false, nil)

		_, res, err := svc.VerifyUrlExpired(context.Background(), code)
		assert.Error(t, err)
		assert.NotNil(t, res)
	})

	t.Run("验证通过", func(t *testing.T) {
		var (
			targetEmail = &model.VerificationEmail{}
			code        = "code"
		)
		mockEmailRepo.EXPECT().VerifyCode(context.Background(), code).Return(targetEmail, true, nil)

		email, res, err := svc.VerifyUrlExpired(context.Background(), code)
		assert.NoError(t, err)
		assert.Nil(t, res)
		assert.Equal(t, targetEmail, email)
	})
}
