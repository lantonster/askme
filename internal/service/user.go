package service

import (
	"context"
	"fmt"
	"time"

	"github.com/lantonster/askme/internal/model"
	"github.com/lantonster/askme/internal/repo"
	"github.com/lantonster/askme/internal/schema"
	"github.com/lantonster/askme/pkg/errors"
	"github.com/lantonster/askme/pkg/handler"
	"github.com/lantonster/askme/pkg/i18n"
	"github.com/lantonster/askme/pkg/log"
	"github.com/lantonster/askme/pkg/reason"
	"github.com/lantonster/askme/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	// GetUserByUserId 根据用户 ID 获取用户信息
	GetUserByUserId(c context.Context, userId int64) (*schema.GetUserByUserIdRes, error)

	// LoginByEmail 邮箱登录
	LoginByEmail(c context.Context, req *schema.LoginByEmailReq) (*schema.LoginByEmailRes, error)

	// RegisterUserByEmail 邮箱注册
	RegisterUserByEmail(c context.Context, req *schema.RegisterUserByEmailReq) (*schema.RegisterUserByEmailRes, []*validator.FieldError, error)

	// VerifyEmail 邮箱验证
	VerifyEmail(c context.Context, req *schema.VerifyEmailReq) (*schema.VerifyEmailRes, error)
}

type UserServiceImpl struct {
	*repo.Repo
}

func NewUserService(repo *repo.Repo) UserService {
	return &UserServiceImpl{Repo: repo}
}

// encryptPassword 密码加密
func (s *UserServiceImpl) encryptPassword(c context.Context, password string) (string, error) {
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.WithContext(c).Errorf("密码 %s 加密失败: %v", password, err)
		return "", fmt.Errorf("密码加密失败: %+w", err)
	}
	return string(hashPwd), nil
}

func (s *UserServiceImpl) GetUserByUserId(c context.Context, userId int64) (*schema.GetUserByUserIdRes, error) {
	user, exist, err := s.UserRepo.GetUserById(c, userId)
	if err != nil {
		log.WithContext(c).Errorf("查询用户 [%d] 信息失败: %v", userId, err)
		return nil, err
	} else if !exist {
		log.WithContext(c).Errorf("用户 [%d] 不存在", userId)
		return nil, errors.BadRequest(reason.UserNotFound)
	} else if user.Status == model.UserStatusDeleted {
		return nil, errors.Unauthorized(reason.UnauthorizedError)
	}

	// TODO role id

	return &schema.GetUserByUserIdRes{
		User:         *user,
		HavePassword: len(user.Password) > 0,
	}, nil
}

func (s *UserServiceImpl) LoginByEmail(c context.Context, req *schema.LoginByEmailReq) (*schema.LoginByEmailRes, error) {
	login, err := siteInfoService.GetSiteLogin(c)
	if err != nil {
		log.WithContext(c).Errorf("获取站点登录信息失败: %v", err)
		return nil, err
	}

	if !login.AllowPasswordLogin {
		return nil, errors.BadRequest(reason.NotAllowedLoginViaPassword)
	}

	user, exist, err := s.UserRepo.GetUserByEmail(c, req.Email)
	if err != nil {
		log.WithContext(c).Errorf("查询用户 [%s] 信息失败: %v", req.Email, err)
		return nil, err
	} else if !exist || user.Status == model.UserStatusDeleted {
		log.WithContext(c).Errorf("用户 [%s] 不存在或已删除", req.Email)
		return nil, errors.BadRequest(reason.EmailOrPasswordWrong)
	}

	// TODO external login

	if err := s.UserRepo.UpdateLastLoginDate(c, user.Id); err != nil {
		log.WithContext(c).Errorf("更新用户 [%d] 最后登录时间失败: %v", user.Id, err)
		return nil, err
	}

	// TODO role id

	res := &schema.LoginByEmailRes{User: *user, HavePassword: len(user.Password) > 0}
	if res.AccessToken, res.VisitToken, err = authService.SetUserCacheInfo(c, user); err != nil {
		log.WithContext(c).Errorf("生成用户相关 token 发生错误: %v", err)
		return nil, err
	}

	if res.RoleId == model.RoleIdAdminID {
		// TODO set admin chache
	}

	return res, nil
}

// RegisterUserByEmail 通过电子邮件注册用户。
//
// 参数:
//   - c: 上下文
//   - req: 包含注册用户所需信息的请求
//
// 返回:
//   - *schema.RegisterUserByEmailRes: 注册响应
//   - []*validator.FieldError: 字段验证错误列表
//   - error: 可能返回的错误
func (s *UserServiceImpl) RegisterUserByEmail(c context.Context, req *schema.RegisterUserByEmailReq) (
	res *schema.RegisterUserByEmailRes, fieldErr []*validator.FieldError, err error) {
	res = &schema.RegisterUserByEmailRes{}

	// 查询指定邮箱是否已存在，如果查询时出错
	if _, exist, err := s.UserRepo.GetUserByEmail(c, req.Email); err != nil {
		log.WithContext(c).Errorf("邮箱注册前查询邮箱是否存在失败: %v", err)
		return nil, nil, err
	} else if exist {
		// 如果邮箱已存在，记录错误日志并返回相应错误
		log.WithContext(c).Errorf("邮箱 [%s] 已存在", req.Email)
		return nil, []*validator.FieldError{{
			Field: "e_mail",
			Error: i18n.Tr(handler.GetLangByCtx(c), reason.UserEmailDuplicate),
		}}, errors.BadRequest(reason.UserEmailDuplicate)
	}

	// 创建新用户对象
	user := &model.User{
		Email:         req.Email,
		DisplayName:   req.Name,
		IpInfo:        req.IP,
		Status:        model.UserStatusAvailable,
		MailStatus:    model.EmailStatusToBeVerified,
		LastLoginDate: time.Now().Unix(),
	}
	// 对用户密码进行加密，如果加密时出错
	if user.Password, err = s.encryptPassword(c, req.Pass); err != nil {
		return nil, nil, err
	}
	// 为用户生成唯一用户名，如果生成时出错
	if user.Username, err = s.UserRepo.GenerateUniqueUsername(c, req.Name); err != nil {
		return nil, []*validator.FieldError{{
			Field: "name",
			Error: i18n.Tr(handler.GetLangByCtx(c), reason.UsernameInvalid),
		}}, err
	}
	// 创建用户，如果创建时出错
	if err = s.UserRepo.CreateUser(c, user); err != nil {
		log.WithContext(c).Errorf("创建用户失败: %v", err)
		return nil, nil, err
	}

	// 异步发送注册验证邮件
	go emailService.SendRegisterVerificationEmail(c, user.Id, user.Email)

	// 设置用户相关 token
	res.User = *user
	res.RoleId = model.RoleIdUser
	res.AccessToken, res.VisitToken, err = authService.SetUserCacheInfo(c, user)
	if err != nil {
		log.WithContext(c).Errorf("生成用户相关 token 发生错误: %v", err)
		return nil, nil, err
	}

	return res, nil, nil
}

// VerifyEmail 函数用于验证用户的电子邮件。
//
// 参数:
//   - c: 上下文
//   - req: 包含要验证的电子邮件的请求
//
// 返回:
//   - *schema.VerifyEmailRes: 验证电子邮件的响应
//   - error: 可能返回的错误
func (s *UserServiceImpl) VerifyEmail(c context.Context, req *schema.VerifyEmailReq) (*schema.VerifyEmailRes, error) {
	res := &schema.VerifyEmailRes{}

	// 通过电子邮件查询用户，如果查询失败则记录错误并返回
	user, exist, err := s.UserRepo.GetUserByEmail(c, req.Email.Email)
	if err != nil {
		log.WithContext(c).Errorf("通过邮箱 [%s] 查询用户失败: %v", req.Email.Email, err)
		return nil, err
	} else if !exist {
		log.WithContext(c).Infof("邮箱 [%s] 未注册", req.Email.Email)
		return nil, errors.BadRequest(reason.UserNotFound)
	}

	// 如果用户的邮件状态为待验证，更新为可用
	if user.MailStatus == model.EmailStatusToBeVerified {
		user.MailStatus = model.EmailStatusAvailable
		if err = s.UserRepo.UpdateEmailStatus(c, user.Id, user.MailStatus); err != nil {
			log.WithContext(c).Errorf("更新用户 [%d] 邮箱状态 [%s] 失败: %v", user.Id, user.MailStatus, err)
			return nil, err
		}
	}

	// 生成用户相关的 token，如果生成出错则记录错误并返回
	if res.AccessToken, res.VisitToken, err = authService.SetUserCacheInfo(c, user); err != nil {
		log.WithContext(c).Errorf("生成用户相关 token 发生错误: %v", err)
		return nil, err
	}

	if err := activityService.ActivateUser(c, user); err != nil {
		log.WithContext(c).Errorf("用户 [%d] 激活创建活动发生错误: %v", user.Id, err)
		return nil, err
	}

	// TODO three-party login

	return res, nil
}
