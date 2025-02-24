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

	// GetAccessTokenByVisitToken 从缓存中获取 visit token 对应的 access token。
	GetAccessTokenByVisitToken(c context.Context, visitToken string) (token string, err error)

	// GetUserCacheByAccessToken 从缓存中获取 access token 对应的用户信息。
	GetUserCacheByAccessToken(c context.Context, accessToken string) (user *model.UserInfo, err error)

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

// GetAccessTokenByVisitToken 从缓存中获取 visit token 对应的 access token。
//
// 参数:
//   - c: 上下文
//   - visitToken: 访问令牌
//
// 返回:
//   - string: 访问令牌，如果获取成功则不为空字符串
//   - error: 可能返回的错误
func (r *AuthRepoImpl) GetAccessTokenByVisitToken(c context.Context, visitToken string) (token string, err error) {
	key := fmt.Sprintf(constant.CacheKeyUserAccessToken, visitToken)

	accessToken, exist, err := r.Cache.GetString(c, key)
	if err != nil {
		log.WithContext(c).Errorf("获取 access token 缓存发生错误: %v", err)
		return "", errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	} else if !exist {
		log.WithContext(c).Infof("access token 缓存不存在: %s", visitToken)
		return "", errors.Unauthorized(reason.UnauthorizedError)
	}

	return accessToken, nil
}

// GetUserCache 函数从缓存中获取指定访问令牌对应的用户信息。
//
// 参数:
//   - c: 上下文
//   - accessToken: 访问令牌
//
// 返回:
//   - *model.UserInfo: 用户信息，如果获取成功则不为 nil
//   - error: 可能返回的错误
func (r *AuthRepoImpl) GetUserCacheByAccessToken(c context.Context, accessToken string) (user *model.UserInfo, err error) {
	user = &model.UserInfo{}

	key := fmt.Sprintf(constant.CacheKeyUserInfo, accessToken)
	if exist, err := r.Cache.GetObj(c, key, user); err != nil {
		log.WithContext(c).Errorf("获取用户信息缓存发生错误: %v", err)
		return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	} else if !exist {
		log.WithContext(c).Infof("用户信息缓存不存在: %s", accessToken)
		return nil, errors.Unauthorized(reason.UnauthorizedError)
	}

	return user, nil
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
