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

func TestFirstSiteInfoByType(t *testing.T) {
	repo := repo.NewSiteInfoRepo(data)

	siteInfo := &model.SiteInfo{
		Id:        1,
		Type:      string(model.SiteInfoTypeBranding),
		Content:   "content",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	bytes, _ := json.Marshal(siteInfo)

	t.Run("from cache", func(t *testing.T) {
		key := fmt.Sprintf(constant.CacheKeySiteInfo, siteInfo.Type)
		rdb.Set(context.Background(), key, string(bytes), 0)

		_siteInfo, err := repo.FirstSiteInfoByType(context.Background(), model.SiteInfoType(siteInfo.Type))
		assert.NoError(t, err)
		assert.Equal(t, siteInfo, _siteInfo)
	})

	t.Run("from db", func(t *testing.T) {
		db.Create(siteInfo)

		_siteInfo, err := repo.FirstSiteInfoByType(context.Background(), model.SiteInfoType(siteInfo.Type))
		assert.NoError(t, err)
		assert.Equal(t, siteInfo, _siteInfo)
	})
}
