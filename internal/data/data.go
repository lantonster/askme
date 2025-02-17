package data

import (
	"context"
	"fmt"

	"github.com/lantonster/askme/internal/conf"
	"github.com/lantonster/askme/pkg/cache"
	"github.com/lantonster/askme/pkg/cache/redis"
	"github.com/lantonster/askme/pkg/log"
	"github.com/lantonster/askme/pkg/orm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Data struct {
	DB    *gorm.DB
	Cache cache.Cache
}

func NewData(db *gorm.DB, cache cache.Cache) *Data {
	orm.SetDefault(db)
	return &Data{
		DB:    db,
		Cache: cache,
	}
}

func NewCache(config *conf.Config) cache.Cache {
	cache := redis.NewCache(config.Redis.Addr, config.Redis.User, config.Redis.Password)
	log.WithContext(context.Background()).Infof("redis [%s] connected", config.Redis.Addr)
	return cache
}

func NewGormDB(config *conf.Config) *gorm.DB {
	// 构建 MySQL 的 DSN (Data Source Name)
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True",
		config.MySQL.User,
		config.MySQL.Password,
		config.MySQL.Addr,
		config.MySQL.DBName,
	)

	// 设置时区为上海
	dsn = dsn + "&loc=Asia%2FShanghai"

	var db *gorm.DB
	var err error

	// 使用 Gorm 开启 MySQL 连接，并配置外键约束迁移行为
	if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   log.NewGormLogger(config.MySQL.Debug),
	}); err != nil {
		panic(err)
	}
	log.WithContext(context.Background()).Infof("db [%s/%s] connected", config.MySQL.Addr, config.MySQL.DBName)

	return db
}
