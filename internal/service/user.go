package service

import (
	"context"
	"fmt"
	"time"

	"github.com/lantonster/askme/internal/model"
	"github.com/lantonster/askme/internal/repo"
	"github.com/lantonster/askme/internal/schema"
	"github.com/lantonster/askme/pkg/errors"
	"github.com/lantonster/askme/pkg/log"
	"github.com/lantonster/askme/pkg/reason"
	"github.com/lantonster/askme/pkg/token"
	"github.com/lantonster/askme/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	// RegisterUserByEmail 通过邮箱注册账号
	RegisterUserByEmail(c context.Context, req *schema.RegisterUserByEmailReq) (*schema.RegisterUserByEmailRes, []*validator.FieldError, error)
}

type UserServiceImpl struct {
	repo *repo.Repo
}

func NewUserService(repo *repo.Repo) UserService {
	return &UserServiceImpl{repo: repo}
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

func (s *UserServiceImpl) getSiteUrl(c context.Context) string {
	general, err := siteInfoService.GetSiteGeneral(c)
	if err != nil {
		log.WithContext(c).Errorf("获取站点网址发生错误: %v", err)
		return ""
	}
	return general.SiteUrl
}

func (s *UserServiceImpl) RegisterUserByEmail(c context.Context, req *schema.RegisterUserByEmailReq) (
	res *schema.RegisterUserByEmailRes, fieldErr []*validator.FieldError, err error) {

	if _, exist, err := s.repo.User.GetUserByEmail(c, req.Email); err != nil {
		log.WithContext(c).Errorf("邮箱注册前查询邮箱是否存在失败: %v", err)
		return nil, nil, err
	} else if exist {
		log.WithContext(c).Errorf("邮箱 [%s] 已存在", req.Email)
		return nil, []*validator.FieldError{{Field: "e_mail", Error: reason.UserEmailDuplicate}}, errors.BadRequest(reason.UserEmailDuplicate)
	}

	user := &model.User{
		Email:         req.Email,
		DisplayName:   req.Name,
		IP:            req.IP,
		Status:        model.UserStatusAvailable,
		MailStatus:    model.EmailStatusToBeVerified,
		LastLoginDate: time.Now().Unix(),
	}
	if user.Password, err = s.encryptPassword(c, req.Pass); err != nil {
		return nil, nil, err
	}
	if user.Username, err = s.repo.User.GenerateUniqueUsername(c, req.Name); err != nil {
		return nil, []*validator.FieldError{{Field: "name", Error: reason.UsernameInvalid}}, err
	}
	if err = s.repo.User.CreateUser(c, user); err != nil {
		log.WithContext(c).Errorf("创建用户失败: %v", err)
		return nil, nil, err
	}

	code := token.GenerateToken()
	url := fmt.Sprintf("%s/users/account-activation?code=%s", s.getSiteUrl(c), code)
	fmt.Printf("url: %v\n", url)

	// TODO role id
	// TODO avatar

	return nil, nil, nil
}
