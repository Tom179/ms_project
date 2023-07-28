package repo

import (
	"context"
	"time"
)

type Cache interface { //存储接口
	Put(ctx context.Context, key, value string, expire time.Duration) error //传入超时上下文
	Get(ctx context.Context, key string) (string, error)
}
