package repo

import (
	"testing"

	redisv2 "github.com/go-redis/redis/v8"
	_data "github.com/lantonster/askme/internal/data"
	"github.com/lantonster/askme/internal/model"
	"github.com/lantonster/askme/pkg/cache/redis"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	data *_data.Data
	db   *gorm.DB
	rdb  *redisv2.Client
)

func TestMain(m *testing.M) {
	cache := redis.NewMemoryCache()
	rdb = cache.Client
	db = NewDB()

	data = _data.NewData(db, cache)

	m.Run()
}

func NewDB() *gorm.DB {
	// 创建内存中的 SQLite 数据库连接
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		model.Activity{},
		model.Config{},
		model.Role{},
		model.SiteInfo{},
		model.User{},
	)

	return db
}
