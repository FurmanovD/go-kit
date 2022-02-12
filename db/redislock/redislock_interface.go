package redislock

import (
	"context"
	"errors"
	"time"
)

type RedisLockError error

var (
	// these errors describes a result of a record lock
	Ok             = RedisLockError(nil)
	AlreadyLocked  = RedisLockError(errors.New("a key is already locked"))
	Uninitialized  = RedisLockError(errors.New("locker or redis client is not initialized"))
	EmptyKey       = RedisLockError(errors.New("key to lock is empty"))
	UnlockRequired = RedisLockError(errors.New("locker is already locked another key"))
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
