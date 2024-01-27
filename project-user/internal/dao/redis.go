package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"test.com/project-user/config"
	"time"
)

var Rc *RedisCache //内部单例
type RedisCache struct {
	rdb *redis.Client
}

func (rc *RedisCache) Put(ctx context.Context, key, value string, expire time.Duration) error {
	err := rc.rdb.Set(ctx, key, value, expire).Err()
	return err
}
func (rc *RedisCache) Get(ctx context.Context, key string) (string, error) {
	result, err := rc.rdb.Get(ctx, key).Result()
	return result, err
}

func init() { //连接redis，初始化(赋值)内部单例，外界直接引用这个单例
	rdb := redis.NewClient(config.C.ReadRedisConfig())
	Rc = &RedisCache{rdb: rdb}
}
