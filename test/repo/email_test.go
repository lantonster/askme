package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/lantonster/askme/internal/constant"
	"github.com/lantonster/askme/internal/model"
	"github.com/lantonster/askme/internal/repo"
	"github.com/stretchr/testify/assert"
)

func TestStoreVerificationEmail(t *testing.T) {
	repo := repo.NewEmailRepo(data)

	code := "code"
	email := &model.VerificationEmail{
		UserId:                   1,
		Email:                    "email",
		BindingKey:               "bindingKey",
		SkipValidationLatestCode: true,
	}
	bytes, _ := json.Marshal(email)

	err := repo.StoreVerificationEmail(context.Background(), email.UserId, code, email, time.Minute)
	assert.NoError(t, err)

	key := fmt.Sprintf(constant.CacheKeyVerificationEmailLatestCode, email.UserId)
	out, _ := rdb.Get(context.Background(), key).Result()
	assert.Equal(t, code, out)

	key = fmt.Sprintf(constant.CacheKeyVerificationEmail, code)
	out, _ = rdb.Get(context.Background(), key).Result()
	assert.Equal(t, string(bytes), out)
}

func TestVerifyCode(t *testing.T) {
	repo := repo.NewEmailRepo(data)

	skipCode := []string{"skip1", "skip2"}
	noSkipCode := []string{"noSkip1", "noSkip2"}
	skip := &model.VerificationEmail{
		UserId:                   1,
		Email:                    "email",
		BindingKey:               "bindingKey",
		SkipValidationLatestCode: true,
	}
	noSkip := &model.VerificationEmail{
		UserId:                   2,
		Email:                    "email",
		BindingKey:               "bindingKey",
		SkipValidationLatestCode: false,
	}
	init := func() {
		for _, code := range skipCode {
			repo.StoreVerificationEmail(context.Background(), skip.UserId, code, skip, time.Minute)
		}
		for _, code := range noSkipCode {
			repo.StoreVerificationEmail(context.Background(), noSkip.UserId, code, noSkip, time.Minute)
		}
	}

	t.Run("code expired", func(t *testing.T) {
		_, success, err := repo.VerifyCode(context.Background(), "expired")
		assert.NoError(t, err)
		assert.False(t, success)
	})

	t.Run("skip validateion latest code", func(t *testing.T) {
		init()

		email, success, err := repo.VerifyCode(context.Background(), skipCode[0])
		assert.NoError(t, err)
		assert.True(t, success)
		assert.Equal(t, skip, email)
	})

	t.Run("no skip", func(t *testing.T) {
		init()

		_, success, err := repo.VerifyCode(context.Background(), noSkipCode[0])
		assert.NoError(t, err)
		assert.False(t, success)

		email, success, err := repo.VerifyCode(context.Background(), noSkipCode[1])
		assert.NoError(t, err)
		assert.True(t, success)
		assert.Equal(t, noSkip, email)
	})
}
