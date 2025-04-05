//go:generate mockery --with-expecter --name=RedisLock --testonly --inpackage --filename=redislock_mock.go
package redislock

import (
	"context"
	"time"
)

// TODO(DF) Add later some "lockerID" parameter to unlock all records locked by the server if it crashed

// RedisLock describes a locker interface.
type RedisLock interface {
	Lock(ctx context.Context, ttl time.Duration) Error
	ObtainLock(
		ctx context.Context,
		ttl time.Duration,
		timeout time.Duration,
		loopPeriod time.Duration,
	) Error
	Unlock(ctx context.Context) Error
}
