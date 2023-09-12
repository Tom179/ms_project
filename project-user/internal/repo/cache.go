package repo //repo包用于定义各种存储接口

import (
	"context"
	"time"
)

type Cache interface { //存储接口
	Put(ctx context.Context, key, value string, expire time.Duration) error //传入超时上下文
	Get(ctx context.Context, key string) (string, error)
}
