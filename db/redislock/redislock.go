package redislock

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	//
	lockValue     = 1
	lockKeyPrefix = "lock-"
)

type clock interface {
	Now() time.Time
}

type redisClient interface {
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

type redisLocker struct {
	mutex   sync.Mutex
	rclient redisClient
	key     string
	clock   clock
	locked  bool
	//TODO(DF) possibly add a lock-count to allow the same locker lock the same key, e.g. to extend a lock TTL
}

func NewRedisLocker(rc redisClient, key string, clk clock) RedisLock {
	return &redisLocker{
		rclient: rc,
		key:     key,
		clock:   clk,
	}
}

// Locks a redis record by creating a key with special name or returns an ErrAlreadyLocked if such a key already exists
func (rl *redisLocker) Lock(ctx context.Context, ttl time.Duration) Error {
	if rl == nil || rl.rclient == nil {
		return ErrUninitialized
	}

	if rl.key == "" {
		return ErrEmptyKey
	}

	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	return rl.tryLock(ctx, ttl)
}

// ObtainLock tries to lock a key until try timeout is reached with loopPeriod pause between tries.
func (rl *redisLocker) ObtainLock(
	ctx context.Context,
	ttl time.Duration,
	timeout time.Duration,
	retryPeriod time.Duration,
) Error {
	if rl == nil || rl.rclient == nil {
		return ErrUninitialized
	}

	if rl.key == "" {
		return ErrEmptyKey
	}

	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := rl.clock.Now()
	timeoutTime := now.Add(timeout)

	// use now variable and !now.After condition to make sure at least one try is done
	// even in case 0 timeout is received
	for ; !now.After(timeoutTime); now = time.Now() {
		tryRes := rl.tryLock(ctx, ttl)
		if tryRes == ErrAlreadyLocked {
			time.Sleep(retryPeriod)
			continue
		}

		return tryRes
	}

	// if we reached this point, it means timeout is reached and another lock still not released
	return ErrAlreadyLocked
}

func (rl *redisLocker) Unlock(ctx context.Context) Error {
	if rl == nil || rl.rclient == nil {
		return ErrUninitialized
	}

	if rl.key == "" {
		return ErrEmptyKey
	}

	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	if rl.locked {
		rl.rclient.Del(ctx, rl.key)
	}
	rl.locked = false

	return nil
}

// tryLock actually locks the record.
// ! No sync.
// ! No parameters validation.
func (rl *redisLocker) tryLock(ctx context.Context, ttl time.Duration) Error {
	lockKey := lockKeyPrefix + rl.key
	if rl.locked {
		return ErrUnlockRequired
	}

	if setRes := rl.rclient.SetNX(ctx, lockKey, lockValue, ttl); setRes.Val() {
		rl.locked = true
		return nil
	}

	return ErrAlreadyLocked
}
