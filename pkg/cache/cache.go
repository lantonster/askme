package cache

import (
	"context"
	"time"
)

type Cache interface {
	Del(c context.Context, key string) error

	// SetObj 函数将对象存储到缓存中。
	SetObj(c context.Context, key string, obj any, ttl time.Duration) (err error)

	// GetObj 函数从缓存中获取对象。
	GetObj(c context.Context, key string, obj any) (exist bool, err error)

	// SetString 将字符串存储到缓存中。
	SetString(c context.Context, key string, value string, ttl time.Duration) (err error)

	// GetString 从缓存中获取字符串。
	GetString(c context.Context, key string) (value string, exist bool, err error)
}
