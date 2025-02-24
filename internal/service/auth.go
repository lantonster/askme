package service

import (
	"context"

	"github.com/lantonster/askme/internal/model"
	"github.com/lantonster/askme/internal/repo"
	"github.com/lantonster/askme/pkg/log"
	"github.com/lantonster/askme/pkg/token"
)

type AuthService interface {
	// CheckVisitToken 检查给定的访问令牌是否有效。
	CheckVisitToken(c context.Context, visitToken string) bool

	// GetUserCacheInfo 获取指定访问令牌对应的用户缓存信息。
	GetUserCacheInfo(c context.Context, accessToken string) (user *model.UserInfo, err error)

	// SetUserCacheInfo 生成用户需要的 token 并缓存用户信息。
	SetUserCacheInfo(c context.Context, user *model.User) (accessToken, visitToken string, err error)
}

type AuthServiceImpl struct {
	*repo.Repo
}

func NewAuthService(repo *repo.Repo) AuthService {
	return &AuthServiceImpl{Repo: repo}
}

// CheckVisitToken 检查给定的访问令牌是否有效。
//
// 参数:
//   - c: 上下文
//   - visitToken: 要检查的访问令牌
//
// 返回:
//   - bool: 如果访问令牌有效（对应的访问令牌不为空）则返回 true，否则返回 false
func (s *AuthServiceImpl) CheckVisitToken(c context.Context, visitToken string) bool {
	accessToken, err := s.AuthRepo.GetAccessTokenByVisitToken(c, visitToken)
	if err != nil {
		log.WithContext(c).Errorf("获取 visitToken [%s] 对应的 accessToken 时发生错误: %v", visitToken, err)
		return false
	}
	return accessToken != ""
}

// GetUserCacheInfo 获取指定访问令牌对应的用户缓存信息。
//
// 参数:
//   - c: 上下文
//   - accessToken: 访问令牌
//
// 返回:
//   - *model.UserInfo: 用户信息，如果获取成功则不为 nil
//   - error: 可能返回的错误
func (s *AuthServiceImpl) GetUserCacheInfo(c context.Context, accessToken string) (user *model.UserInfo, err error) {
	user, err = s.AuthRepo.GetUserCacheByAccessToken(c, accessToken)
	if err != nil {
		log.WithContext(c).Errorf("获取用户 [%d] 信息时发生错误: %v", user.UserId, err)
		return
	}

	return
}

// SetUserCacheInfo 生成用户需要的 token 并缓存用户信息。
//
// 参数:
//   - c: 上下文
//   - user: 用户信息
//
// 返回:
//   - accessToken: 生成的访问令牌
//   - visitToken: 生成的访问令牌
//   - err: 可能返回的错误
func (s *AuthServiceImpl) SetUserCacheInfo(c context.Context, user *model.User) (accessToken, visitToken string, err error) {
	// TODO get user role

	// if role id == admin id { set admin cache }

	info := &model.UserInfo{
		UserId:      user.Id,
		RoleId:      model.RoleIdUser,
		UserStatus:  user.Status,
		EmailStatus: user.MailStatus,
	}

	// 生成访问令牌和访问令牌
	accessToken, visitToken = token.GenerateToken(), token.GenerateToken()
	info.VisitToken = visitToken

	// 设置用户的缓存
	if err := s.AuthRepo.SetUserCache(c, accessToken, info); err != nil {
		log.WithContext(c).Errorf("缓存用户 [%d] 信息时发生错误: %v", info.UserId, err)
		return "", "", err
	}
	return
}
