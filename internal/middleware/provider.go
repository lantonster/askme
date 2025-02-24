package middleware

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewMiddleware,
	NewAuthMiddleware,
	NewAvatarMiddleware,
)

type Middleware struct {
	Auth   *AuthMiddleware
	Avatar *AvatarMiddleware
}

func NewMiddleware(
	auth *AuthMiddleware,
	avatar *AvatarMiddleware,
) *Middleware {
	return &Middleware{
		Auth:   auth,
		Avatar: avatar,
	}
}
