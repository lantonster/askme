package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/lantonster/askme/internal/constant"
	"github.com/lantonster/askme/internal/model"
	"github.com/lantonster/askme/internal/repo"
	"github.com/stretchr/testify/assert"
)

func TestAddUserAccessToken(t *testing.T) {
	repo := repo.NewAuthRepo(data)
	var (
		userId int64  = 1
		token  string = "token"
	)

	err := repo.AddUserAccessToken(context.Background(), userId, token)
	assert.NoError(t, err)

	key := fmt.Sprintf(constant.CacheKeyUserAccessTokenMapping, userId)
	out, _ := rdb.Get(context.Background(), key).Result()
	assert.Equal(t, `{"token":true}`, out)
}

func TestSetUserCache(t *testing.T) {
	repo := repo.NewAuthRepo(data)
	var (
		accessToken string          = "access_token"
		user        *model.UserInfo = &model.UserInfo{
			UserId:     1,
			VisitToken: "visit_token",
		}
	)
	bytes, _ := json.Marshal(user)

	err := repo.SetUserCache(context.Background(), accessToken, user)
	assert.NoError(t, err)

	key := fmt.Sprintf(constant.CacheKeyUserInfo, accessToken)
	out, _ := rdb.Get(context.Background(), key).Result()
	assert.Equal(t, string(bytes), out)

	key = fmt.Sprintf(constant.CacheKeyUserAccessToken, user.VisitToken)
	out, _ = rdb.Get(context.Background(), key).Result()
	assert.Equal(t, accessToken, out)
}
