package schema

import (
	"github.com/lantonster/askme/internal/model"
	"github.com/lantonster/askme/pkg/checker"
	"github.com/lantonster/askme/pkg/validator"
)

type RegisterUserByEmailReq struct {
	Name  string `validate:"required,gte=2,lte=30" json:"name"`          // 用户名
	Email string `validate:"required,email,gt=0,lte=500" json:"e_mail" ` // 邮箱
	Pass  string `validate:"required,gte=8,lte=32" json:"pass"`          // 密码
	IP    string `json:"-" `                                             // IP 地址
}

// Check 方法用于检查 `RegisterUserByEmailReq` 中密码的有效性。
//
// 返回:
//   - []*validator.FieldError: 包含错误信息的字段数组，如果密码不符合要求则包含一个关于密码的错误
//   - error: 始终为 nil
func (r *RegisterUserByEmailReq) Check() (fields []*validator.FieldError, err error) {
	// 检查密码的强度
	if err = checker.CheckPassword(r.Pass, checker.PasswordStrengthHig); err != nil {
		// 如果密码不符合要求，返回一个包含密码错误信息的字段
		return []*validator.FieldError{{
			Field: "pass",
			Error: err.Error(),
		}}, nil
	}
	// 如果密码符合要求，返回空的字段数组和 nil 错误
	return nil, nil
}

type RegisterUserByEmailRes = UserLoginRes

type VerifyEmailReq struct {
	Code  string                   `validate:"required,gt=0,lte=500" form:"code"` // 验证码
	Email *model.VerificationEmail `json:"-"`                                     // 验证邮箱内容
}

type VerifyEmailRes = UserLoginRes

type UserLoginRes struct {
	model.User

	RoleId      int64  `json:"role_id"`      // 角色 ID
	AccessToken string `json:"access_token"` // 访问令牌
	VisitToken  string `json:"visit_token"`  // 访问令牌
}
