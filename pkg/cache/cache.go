package cache

import (
	"context"
	"time"
)

type Cache interface {
	// SetObj 函数将对象存储到缓存中。
	//
	// 参数:
	//   - c: 上下文
	//   - key: 缓存键
	//   - obj: 要存储的对象
	//   - ttl: 过期时间
	// 返回:
	//   - err: 可能出现的错误
	SetObj(c context.Context, key string, obj any, ttl time.Duration) (err error)

	// GetObj 函数从缓存中获取对象。
	//
	// 参数:
	//   - c: 上下文
	//   - key: 缓存键
	//   - obj: 用于存储获取到的数据的对象指针，必须为指针才能正确获取到值
	// 返回:
	//   - exist: 是否存在对应键的值
	//   - err: 可能出现的错误
	GetObj(c context.Context, key string, obj any) (exist bool, err error)
}
