package service

import (
	"context"

	"github.com/lantonster/askme/internal/model"
	"github.com/lantonster/askme/internal/repo"
	"github.com/lantonster/askme/pkg/log"
	"github.com/lantonster/askme/pkg/token"
)

type AuthService interface {
	// SetUserCacheInfo 生成用户需要的 token 并缓存用户信息。
	SetUserCacheInfo(c context.Context, user *model.User) (accessToken, visitToken string, err error)
}

type AuthServiceImpl struct {
	*repo.Repo
}

func NewAuthService(repo *repo.Repo) AuthService {
	return &AuthServiceImpl{Repo: repo}
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
