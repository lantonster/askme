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

	// VerifyCode 验证邮箱验证码。
	VerifyCode(c context.Context, code string) (email *model.VerificationEmail, success bool, err error)
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
		log.WithContext(c).Errorf("缓存用户 [%d] 最新邮件验证码失败: %v", userId, err)
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	// 缓存邮件内容
	key = fmt.Sprintf(constant.CacheKeyVerificationEmail, code)
	if err := r.Cache.SetObj(c, key, email, ttl); err != nil {
		log.WithContext(c).Errorf("缓存验证邮箱 [%s] 内容失败: %v", code, err)
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

// VerifyCode 验证邮箱验证码。
//
// 参数:
//   - c: 上下文
//   - code: 要验证的验证码
//
// 返回:
//   - *model.VerificationEmail: 验证邮件的相关信息，如果验证成功则不为 nil
//   - bool: 验证是否成功
//   - error: 可能返回的错误
func (r *EmailRepoImpl) VerifyCode(c context.Context, code string) (email *model.VerificationEmail, success bool, err error) {
	email = &model.VerificationEmail{}

	// 根据验证码生成缓存键，尝试从缓存获取对应的验证邮件信息
	key := fmt.Sprintf(constant.CacheKeyVerificationEmail, code)
	if exist, err := r.Cache.GetObj(c, key, email); err != nil {
		log.WithContext(c).Errorf("获取验证邮箱 [%s] 内容失败: %v", code, err)
		return nil, false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	} else if !exist {
		log.WithContext(c).Infof("验证邮箱 [%s] 内容已过期", code)
		return nil, false, nil
	}

	// 删除已验证的缓存内容
	if err := r.Cache.Del(c, key); err != nil {
		log.WithContext(c).Errorf("删除验证邮箱 [%s] 内容失败: %v", code, err)
	}

	// 如果设置了跳过验证最新验证码，直接返回成功
	if email.SkipValidationLatestCode {
		return email, true, nil
	}

	// 根据用户 ID 生成最新验证码的缓存键
	key = fmt.Sprintf(constant.CacheKeyVerificationEmailLatestCode, email.UserId)
	// 获取用户最新验证码
	latestCode, exist, err := r.Cache.GetString(c, key)
	if err != nil {
		log.WithContext(c).Errorf("获取用户 [%d] 最新邮件验证码失败: %v", email.UserId, err)
		return nil, false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	} else if !exist {
		log.WithContext(c).Infof("用户 [%d] 最新邮件验证码已过期", email.UserId)
		return nil, false, nil
	}

	// 比较输入的验证码与获取的最新验证码是否匹配
	if code != latestCode {
		log.WithContext(c).Infof("邮箱验证码 [%s] 与用户 [%d] 最新邮件验证码 [%s] 不匹配", code, email.UserId, latestCode)
		return nil, false, nil
	}

	return email, true, nil
}
