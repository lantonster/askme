package service

import (
	"context"
	"fmt"
	"mime"

	"github.com/lantonster/askme/internal/constant"
	"github.com/lantonster/askme/internal/model"
	"github.com/lantonster/askme/internal/repo"
	"github.com/lantonster/askme/internal/schema"
	"github.com/lantonster/askme/pkg/errors"
	"github.com/lantonster/askme/pkg/handler"
	"github.com/lantonster/askme/pkg/i18n"
	"github.com/lantonster/askme/pkg/log"
	"github.com/lantonster/askme/pkg/reason"
	"github.com/lantonster/askme/pkg/token"
	"gopkg.in/gomail.v2"
)

type EmailService interface {
	// Send 发送邮件。
	Send(c context.Context, email, subject, body string) error

	// SendRegisterVerificationEmail 发送注册验证邮件
	SendRegisterVerificationEmail(c context.Context, userId int64, email string) error

	// VerifyUrlExpired 验证 URL 是否过期。
	VerifyUrlExpired(c context.Context, code string) (*model.VerificationEmail, *schema.ForbiddenRes, error)
}

type EmailServiceImpl struct {
	*repo.Repo
}

func NewEmailService(repo *repo.Repo) EmailService {
	return &EmailServiceImpl{Repo: repo}
}

func (s *EmailServiceImpl) getSiteUrl(c context.Context) string {
	general, err := siteInfoService.GetSiteGeneral(c)
	if err != nil {
		return ""
	}
	return general.SiteUrl
}

// Send 发送邮件。
//
// 参数:
//   - c: 上下文
//   - email: 收件人邮箱地址
//   - subject: 邮件主题
//   - body: 邮件正文
//
// 返回: 可能返回的错误
func (s *EmailServiceImpl) Send(c context.Context, email, subject, body string) error {
	// 记录尝试发送邮件的日志
	log.WithContext(c).Infof("尝试发送邮件到 %s", email)

	// 获取邮件配置
	config, err := configService.GetConfigEmail(c)
	if err != nil {
		log.WithContext(c).Errorf("获取邮件配置失败: %v", err)
		return err
	}

	// 如果未配置 SMTP 服务器
	if len(config.SMTPHost) == 0 {
		log.WithContext(c).Errorf("未配置 SMTP 服务器")
		return errors.InternalServer(reason.DatabaseError).WithMsg("未配置 SMTP 服务器").WithStack()
	}

	// 创建新的邮件消息
	m := gomail.NewMessage()
	// 设置发件人信息
	fromName := mime.QEncoding.Encode("utf-8", config.FromName)
	m.SetHeader("From", fmt.Sprintf("%s <%s>", fromName, config.FromEmail))
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// 创建邮件发送拨号器
	d := gomail.NewDialer(config.SMTPHost, config.SMTPPort, config.SMTPUsername, config.SMTPPassword)
	// 根据配置设置是否使用 SSL 或 TLS
	d.SSL = config.IsSSL()

	// 发送邮件，如果失败记录错误日志并返回错误
	if err := d.DialAndSend(m); err != nil {
		log.WithContext(c).Errorf("发送邮件到 %s 失败: %v", email, err)
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	log.WithContext(c).Infof("发送邮件到 %s 成功", email)
	return nil
}

// SendRegisterVerificationEmail 发送注册验证邮件。
// 参数:
//   - c: 上下文
//   - userId: 用户 ID
//   - email: 收件人邮箱地址
//
// 返回: 可能返回的错误
func (s *EmailServiceImpl) SendRegisterVerificationEmail(c context.Context, userId int64, email string) error {
	// 生成验证码
	code := token.GenerateToken()
	// 生成注册 URL
	registerUrl := fmt.Sprintf("%s/users/account-activation?code=%s", s.getSiteUrl(c), code)

	// 获取站点通用信息，如果获取失败则返回错误
	general, err := siteInfoService.GetSiteGeneral(c)
	if err != nil {
		return err
	}

	// 获取语言并初始化注册模板数据
	lang := handler.GetLangByCtx(c)
	template := &schema.RegisterTemplateData{SiteName: general.Name, RegisterUrl: registerUrl}

	// 翻译并获取注册邮件标题和正文
	title := i18n.TrWithData(lang, constant.EmailTplKeyRegisterTitle, template)
	body := i18n.TrWithData(lang, constant.EmailTplKeyRegisterBody, template)
	verification := &model.VerificationEmail{Email: email, UserId: userId}

	// 存储邮件验证码，如果存储失败则记录错误并返回
	if err := s.EmailRepo.StoreVerificationEmail(c, userId, code, verification, constant.CacheTimeVerificationEmail); err != nil {
		log.WithContext(c).Errorf("缓存邮件验证码失败: %v", err)
		return err
	}

	// 发送邮件，如果发送失败则记录错误并返回
	if err := s.Send(c, email, title, body); err != nil {
		log.WithContext(c).Errorf("发送邮件失败: %v", err)
		return err
	}

	return nil
}

// VerifyUrlExpired 验证 URL 是否过期。
//
// 参数:
//   - c: 上下文
//   - code: 验证码
//
// 返回:
//   - *model.VerificationEmail: 验证邮件信息，如果验证成功则不为 nil
//   - *schema.ForbiddenRes: 禁止响应，如果验证码过期则不为 nil
//   - error: 可能返回的错误
func (s *EmailServiceImpl) VerifyUrlExpired(c context.Context, code string) (*model.VerificationEmail, *schema.ForbiddenRes, error) {
	// 调用 EmailRepo 的 VerifyCode 方法进行验证码验证
	email, success, err := s.EmailRepo.VerifyCode(c, code)
	if err != nil {
		log.WithContext(c).Errorf("验证邮件验证码失败: %v", err)
		return nil, nil, err
	} else if !success {
		log.WithContext(c).Infof("邮件验证码 [%s] 已过期", code)
		return nil, &schema.ForbiddenRes{Type: schema.ForbiddenReasonTypeURLExpired}, errors.Forbidden(reason.EmailVerifyURLExpired)
	}

	return email, nil, nil
}
