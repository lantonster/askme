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

func TestFirstConfigByKey(t *testing.T) {
	repo := repo.NewConfigRepo(data)

	t.Run("从缓存中获得", func(t *testing.T) {
		config := &model.Config{Key: "from_cache", Value: "test"}
		cacheKey := fmt.Sprintf(constant.CacheKeyConfig, config.Key)
		bytes, _ := json.Marshal(config)
		rdb.Set(context.Background(), cacheKey, string(bytes), time.Minute)

		_config, err := repo.FirstConfigByKey(context.Background(), config.Key)
		assert.NoError(t, err)
		assert.Equal(t, config, _config)
	})

	t.Run("从数据库中获得", func(t *testing.T) {
		config := &model.Config{Key: "from_db", Value: "test"}
		db.Create(config)

		_config, err := repo.FirstConfigByKey(context.Background(), config.Key)
		assert.NoError(t, err)
		assert.Equal(t, config, _config)

		cacheKey := fmt.Sprintf(constant.CacheKeyConfig, config.Key)
		out, _ := rdb.Get(context.Background(), cacheKey).Result()
		assert.Equal(t, `{"Id":1,"Key":"from_db","Value":"test","Email":null,"UserActivated":0}`, out)
	})
}
