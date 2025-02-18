package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/lantonster/askme/internal/constant"
	"github.com/lantonster/askme/internal/data"
	"github.com/lantonster/askme/internal/model"
	"github.com/lantonster/askme/pkg/errors"
	"github.com/lantonster/askme/pkg/log"
	"github.com/lantonster/askme/pkg/reason"
)

type EmailRepo interface {
	// StoreVerificationEmail 缓存储用户的最新验证码和邮件内容。
	StoreVerificationEmail(c context.Context, userId int64, code string, content *model.VerificationEmail, ttl time.Duration) error
}

type EmailRepoImpl struct {
	*data.Data
}

func NewEmailRepo(data *data.Data) EmailRepo {
	return &EmailRepoImpl{Data: data}
}

// StoreVerificationEmail 缓存储用户的最新验证码和邮件内容。
//
// 参数:
//   - c: 上下文
//   - userId: 用户 ID
//   - code: 验证码
//   - email: 验证邮件的详细信息
//   - ttl: 缓存的生存时间
//
// 返回: 可能返回的错误
func (r *EmailRepoImpl) StoreVerificationEmail(c context.Context, userId int64, code string, email *model.VerificationEmail, ttl time.Duration) error {
	// 缓存最新邮件验证码
	key := fmt.Sprintf(constant.CacheKeyVerificationEmailLatestCode, userId)
	if err := r.Cache.SetString(c, key, code, ttl); err != nil {
		log.WithContext(c).Errorf("缓存用户 %d 最新邮件验证码失败: %v", userId, err)
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	// 缓存邮件内容
	key = fmt.Sprintf(constant.CacheKeyVerificationEmail, code)
	if err := r.Cache.SetObj(c, key, email, ttl); err != nil {
		log.WithContext(c).Errorf("缓存验证邮箱 %s 内容失败: %v", code, err)
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}
