package redislock

import (
	"context"
	"time"
)

// TODO(DF) Add later some "lockerID" parameter to unlock all records locked by the server if it crashed

// RedisLock describes a locker interface.
type RedisLock interface {
	Lock(ctx context.Context, key string, ttl time.Duration) RedisLockError
	ObtainLock(
		ctx context.Context,
		key string,
		ttl time.Duration,
		timeout time.Duration,
		loopPeriod time.Duration,
	) RedisLockError
	Unlock() RedisLockError
}
