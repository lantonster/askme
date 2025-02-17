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

func (r *Cache) SetObj(c context.Context, key string, obj any, ttl time.Duration) (err error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("marshal obj: %w", err)
	}
	return r.Client.Set(c, key, string(bytes), ttl).Err()
}

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
