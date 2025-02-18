package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/lantonster/askme/pkg/cache"
)

var _ cache.Cache = (*Cache)(nil)

type Cache struct {
	Client *redis.Client
}

func NewCache(conn, username, password string) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     conn,
		Username: username,
		Password: password,
	})
	return &Cache{Client: client}
}

// SetObj 函数将对象存储到缓存中。
//
// 参数:
//   - c: 上下文
//   - key: 缓存键
//   - obj: 要存储的对象
//   - ttl: 过期时间
//
// 返回:
//   - err: 可能出现的错误
func (r *Cache) SetObj(c context.Context, key string, obj any, ttl time.Duration) (err error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("marshal obj: %w", err)
	}
	return r.Client.Set(c, key, string(bytes), ttl).Err()
}

// GetObj 函数从缓存中获取对象。
//
// 参数:
//   - c: 上下文
//   - key: 缓存键
//   - obj: 用于存储获取到的数据的对象指针，必须为指针才能正确获取到值
//
// 返回:
//   - exist: 是否存在对应键的值
//   - err: 可能出现的错误
func (r *Cache) GetObj(c context.Context, key string, obj any) (exist bool, err error) {
	data, err := r.Client.Get(c, key).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("get result: %w", err)
	}

	if err = json.Unmarshal([]byte(data), obj); err != nil {
		return false, fmt.Errorf("unmarshal obj: %w", err)
	}
	return true, nil
}

func (r *Cache) SetString(c context.Context, key string, value string, ttl time.Duration) (err error) {
	return r.Client.Set(c, key, value, ttl).Err()
}

func (r *Cache) GetString(c context.Context, key string) (value string, exist bool, err error) {
	data, err := r.Client.Get(c, key).Result()
	if err == redis.Nil {
		return "", false, nil
	} else if err != nil {
		return "", false, fmt.Errorf("get result: %w", err)
	}
	return data, true, nil
}
