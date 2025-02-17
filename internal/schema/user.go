package schema

import (
	"github.com/lantonster/askme/internal/model"
	"github.com/lantonster/askme/pkg/checker"
	"github.com/lantonster/askme/pkg/validator"
)

type RegisterUserByEmailReq struct {
	Name        string `validate:"required,gte=2,lte=30" json:"name"`
	Email       string `validate:"required,email,gt=0,lte=500" json:"e_mail" `
	Pass        string `validate:"required,gte=8,lte=32" json:"pass"`
	CaptchaID   string `json:"captcha_id"`
	CaptchaCode string `json:"captcha_code"`
	IP          string `json:"-" `
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

type RegisterUserByEmailRes = model.User
