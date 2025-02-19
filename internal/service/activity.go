package service

import (
	"context"

	"github.com/lantonster/askme/internal/model"
	"github.com/lantonster/askme/internal/repo"
	"github.com/lantonster/askme/pkg/log"
)

type ActivityService interface {
	// ActivateUser 创建激活用户的相关活动。
	ActivateUser(c context.Context, user *model.User) error
}

type ActivityServiceImpl struct {
	*repo.Repo
}

func NewActivityService(repo *repo.Repo) ActivityService {
	return &ActivityServiceImpl{Repo: repo}
}

// ActivateUser 创建激活用户的相关活动。
//
// 参数:
//   - c: 上下文
//   - user: 用户对象
//
// 返回: 可能返回的错误，如果没有错误则返回 nil
func (s *ActivityServiceImpl) ActivateUser(c context.Context, user *model.User) error {
	// 获取用户激活配置，如果获取失败则记录错误并返回
	config, err := s.ConfigRepo.FirstConfigByKey(c, model.ConfigKeyUserActivated)
	if err != nil {
		log.WithContext(c).Errorf("获取用户激活配置失败: %v", err)
		return err
	}

	// 获取用户激活活动，如果获取失败则记录错误并返回，若存在则直接返回
	if _, exist, err := s.ActivityRepo.GetActivityByType(c, model.ActivityTypeUserActivated); err != nil {
		log.WithContext(c).Errorf("获取用户激活活动失败: %v", err)
		return err
	} else if exist {
		return nil
	}

	// 创建用户激活活动对象
	activity := &model.Activity{
		UserId: user.Id,
		Rank:   config.UserActivated,
		Type:   model.ActivityTypeUserActivated,
	}

	// 增加用户的积分，如果增加失败则记录错误并返回
	if err := s.UserRepo.IncrRank(c, user.Id, user.Rank, activity.Rank); err != nil {
		log.WithContext(c).Errorf("用户 [%d] 激活增加积分失败: %v", user.Id, err)
		return err
	}

	// 创建用户激活活动，如果创建失败则记录错误并返回
	if err := s.ActivityRepo.CreateActivity(c, activity); err != nil {
		log.WithContext(c).Errorf("用户 [%d] 激活创建活动失败: %v", user.Id, err)
		return err
	}

	return nil
}
