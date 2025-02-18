package repo

import (
	"context"
	"fmt"

	"github.com/lantonster/askme/internal/constant"
	"github.com/lantonster/askme/internal/data"
	"github.com/lantonster/askme/internal/model"
	"github.com/lantonster/askme/pkg/errors"
	"github.com/lantonster/askme/pkg/log"
	"github.com/lantonster/askme/pkg/orm"
	"github.com/lantonster/askme/pkg/reason"
	"gorm.io/gorm"
)

type ConfigRepo interface {
	// FirstConfigByKey 根据给定的键在数据库中查找配置信息，并优先从缓存中获取。
	FirstConfigByKey(c context.Context, key string) (config *model.Config, err error)
}

type ConfigRepoImpl struct {
	*data.Data
}

func NewConfigRepo(data *data.Data) ConfigRepo {
	return &ConfigRepoImpl{Data: data}
}

// FirstConfigByKey 根据给定的键在数据库中查找配置信息，并优先从缓存中获取。
//
// 参数:
//   - c: 上下文
//   - key: 配置信息的键
//
// 返回:
//   - *model.Config: 配置对象，如果未找到则为 nil
//   - error: 可能出现的错误
func (r *ConfigRepoImpl) FirstConfigByKey(c context.Context, key string) (config *model.Config, err error) {
	config = &model.Config{}

	cacheKey := fmt.Sprintf(constant.CacheKeyConfig, key)
	// 尝试从缓存中获取配置信息，如果缓存中存在则直接返回
	if exist, _ := r.Cache.GetObj(c, cacheKey, config); exist {
		return config, nil
	}

	// 从数据库中获取配置信息
	if config, err = orm.Q.Config.WithContext(c).Where(orm.Q.Config.Key.Eq(key)).First(); err != nil {
		// 未找到记录
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.InternalServer(reason.ConfigNotFound).WithError(err).WithStack()
		}
		// 其他错误
		return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	// 将配置信息存入缓存
	if err := r.Cache.SetObj(c, cacheKey, config, constant.CacheTimeConfig); err != nil {
		log.WithContext(c).Errorf("缓存配置信息 [%s] 失败: %v", cacheKey, err)
	}
	return config, nil
}
