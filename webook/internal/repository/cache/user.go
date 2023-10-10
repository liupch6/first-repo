package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"geektime/webook/internal/domain"
)

var ErrKeyNotExist = redis.Nil

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, u domain.User) error
	Delete(ctx context.Context, id int64) error
}

type RedisUserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewUserCache(client redis.Cmdable) UserCache {
	return &RedisUserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

func (cache *RedisUserCache) Delete(ctx context.Context, id int64) error {
	return cache.client.Del(ctx, cache.key(id)).Err()
}

func (cache *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	val, err := cache.client.Get(ctx, cache.key(id)).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal(val, &u)
	return u, err
}

func (cache *RedisUserCache) Set(ctx context.Context, u domain.User) error {
	data, err := json.Marshal(u)
	if err != nil {
		return err
	}
	return cache.client.Set(ctx, cache.key(u.Id), data, cache.expiration).Err()
}

func (cache *RedisUserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
