package repo

import (
	"context"
	"fmt"

	"github.com/lantonster/askme/internal/constant"
	"github.com/lantonster/askme/internal/data"
	"github.com/lantonster/askme/internal/model"
	"github.com/lantonster/askme/pkg/errors"
	"github.com/lantonster/askme/pkg/log"
	"github.com/lantonster/askme/pkg/reason"
)

type AuthRepo interface {
	// AddUserAccessToken 为用户添加访问令牌到缓存中的映射。
	AddUserAccessToken(c context.Context, userId int64, accessToken string) error

	// SetUserCache 设置用户相关的缓存信息。
	SetUserCache(c context.Context, accessToken string, user *model.UserInfo) error
}

type AuthRepoImpl struct {
	*data.Data
}

func NewAuthRepo(data *data.Data) AuthRepo {
	return &AuthRepoImpl{Data: data}
}

// AddUserAccessToken 为用户添加访问令牌到缓存中的映射。
//
// 参数:
//   - c: 上下文
//   - userId: 用户 ID
//   - accessToken: 访问令牌
//
// 返回: 可能返回的错误
func (r *AuthRepoImpl) AddUserAccessToken(c context.Context, userId int64, accessToken string) error {
	// 创建一个映射
	mapping := make(map[string]bool)

	// 尝试从缓存中获取用户的访问令牌映射
	key := fmt.Sprintf(constant.CacheKeyUserAccessTokenMapping, userId)
	if _, err := r.Cache.GetObj(c, key, &mapping); err != nil {
		log.WithContext(c).Errorf("获取用户 [%d] access token 映射表发生错误: %v", userId, err)
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	// 将更新后的映射设置回缓存，如果设置时出错
	mapping[accessToken] = true
	if err := r.Cache.SetObj(c, key, &mapping, constant.CacheTimeUserAccessTokenMapping); err != nil {
		log.WithContext(c).Errorf("缓存用户 [%d] access token 映射表发生错误: %v", userId, err)
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

// SetUserCache 设置用户相关的缓存信息。
//
// 参数:
//   - c: 上下文
//   - accessToken: 访问令牌
//   - user: 用户信息
//
// 返回: 可能返回的错误
func (r *AuthRepoImpl) SetUserCache(c context.Context, accessToken string, user *model.UserInfo) error {

	// 将用户信息设置到缓存中
	key := fmt.Sprintf(constant.CacheKeyUserInfo, accessToken)
	if err := r.Cache.SetObj(c, key, user, constant.CacheTimeUserInfo); err != nil {
		log.WithContext(c).Errorf("缓存用户信息 [%+v] 时发生错误: %v", user, err)
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	// 为用户添加访问令牌映射
	if err := r.AddUserAccessToken(c, user.UserId, accessToken); err != nil {
		log.WithContext(c).Errorf("缓存用户 [%d] access token 映射表发生错误: %v", user.UserId, err)
	}

	// 如果用户有访问令牌，将其与访问令牌进行关联并缓存
	if user.VisitToken != "" {
		key = fmt.Sprintf(constant.CacheKeyUserAccessToken, user.VisitToken)
		if err := r.Cache.SetString(c, key, accessToken, constant.CacheTimeUserAccessToken); err != nil {
			log.WithContext(c).Errorf("缓存 access token 发送错误: %v", err)
		}
	}

	return nil
}
