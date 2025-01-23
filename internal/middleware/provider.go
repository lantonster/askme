package middleware

import "github.com/google/wire"

var ProviderSetMiddleware = wire.NewSet(NewAvatarMiddleware)
