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

type SiteInfoRepo interface {
	// FirstSiteInfoByType 根据给定的类型从数据库或缓存中获取第一个站点信息。
	FirstSiteInfoByType(c context.Context, typ model.SiteInfoType) (siteInfo *model.SiteInfo, err error)
}

type siteInfoRepo struct {
	data *data.Data
}

func NewSiteInfoRepo(data *data.Data) SiteInfoRepo {
	return &siteInfoRepo{
		data: data,
	}
}

// FirstSiteInfoByType 根据给定的类型从数据库或缓存中获取第一个站点信息。
//
// 参数:
//   - c: 上下文
//   - typ: 站点信息类型
//
// 返回:
//   - *model.SiteInfo: 站点信息对象指针，如果未找到则为 nil
//   - error: 可能出现的错误
func (r *siteInfoRepo) FirstSiteInfoByType(c context.Context, typ model.SiteInfoType) (siteInfo *model.SiteInfo, err error) {
	// 初始化站点信息对象
	siteInfo = &model.SiteInfo{}

	key := fmt.Sprintf(constant.CacheKeySiteInfo, typ)
	// 尝试从缓存中获取站点信息，如果在缓存中获取成功，直接返回
	if exist, _ := r.data.Cache.GetObj(c, key, siteInfo); exist {
		return siteInfo, nil
	}

	// 从数据库中查询站点信息
	if siteInfo, err = orm.Q.SiteInfo.WithContext(c).Where(orm.Q.SiteInfo.Type.Eq(string(typ))).First(); err != nil {
		// 未找到记录
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.InternalServer(reason.SiteInfoNotFound).WithError(err).WithStack()
		}
		// 其他错误
		return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	// 将站点信息设置到缓存中，发生错误只记录不返回
	if err := r.data.Cache.SetObj(c, key, siteInfo, constant.CacheTimeSiteInfo); err != nil {
		log.WithContext(c).Errorf("cache set obj [%s] error: %s", key, err)
	}
	return siteInfo, nil
}
