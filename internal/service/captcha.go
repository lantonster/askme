package service

type CaptchaService interface {
}

type captchaService struct {
}

func NewCaptchaService() CaptchaService {
	return &captchaService{}
}
