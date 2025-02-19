package repo

import (
	"context"

	"github.com/lantonster/askme/internal/data"
	"github.com/lantonster/askme/internal/model"
	"github.com/lantonster/askme/pkg/errors"
	"github.com/lantonster/askme/pkg/orm"
	"github.com/lantonster/askme/pkg/reason"
)

type ActivityRepo interface {
	// CreateActivity 创建活动
	CreateActivity(c context.Context, activity *model.Activity) error

	// GetActivityByType 根据活动类型获取活动
	GetActivityByType(c context.Context, typ string) (*model.Activity, bool, error)
}

type ActivityRepoImpl struct {
	*data.Data
}

func NewActivityRepo(data *data.Data) ActivityRepo {
	return &ActivityRepoImpl{Data: data}
}

// CreateActivity 创建新的活动。
//
// 参数:
//   - c: 上下文
//   - activity: 要创建的活动对象
//
// 返回: 可能返回的错误，如果没有错误则返回 nil
func (r *ActivityRepoImpl) CreateActivity(c context.Context, activity *model.Activity) error {
	if err := orm.Q.Activity.WithContext(c).Create(activity); err != nil {
		return errors.InternalServer(reason.DatabaseError)
	}
	return nil
}

// GetActivityByType 根据给定的活动类型从数据库中获取活动。
//
// 参数:
//   - c: 上下文
//   - typ: 活动类型
//
// 返回:
//   - *model.Activity: 活动对象，如果找到则不为 nil
//   - bool: 是否找到活动
//   - error: 可能返回的错误
func (r *ActivityRepoImpl) GetActivityByType(c context.Context, typ string) (*model.Activity, bool, error) {
	activities, err := orm.Q.Activity.WithContext(c).Where(orm.Q.Activity.Type.Eq(typ)).Find()
	if err != nil {
		return nil, false, errors.InternalServer(reason.DatabaseError)
	}

	if len(activities) == 0 {
		return nil, false, nil
	}
	return activities[0], true, nil
}
