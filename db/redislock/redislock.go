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

// lockInfo represents an info to unlock.
type lockInfo struct {
	Key string
	Ctx context.Context
}

type redisLocker struct {
	mutex   sync.Mutex
	rclient *redis.Client
	locked  *lockInfo // saves all info to unlock. != nil means locked state
	//TODO(DF) possibly add a lock-count to allow the same locker lock the same key, e.g. to extend a lock TTL
}

func NewRedisLocker(c *redis.Client) RedisLock {
	return &redisLocker{
		rclient: c,
	}
}

// Locks a redis record by creating a key with special name or returns an ErrAlreadyLocked if such a key already exists
func (rl *redisLocker) Lock(ctx context.Context, key string, ttl time.Duration) Error {

	if rl == nil || rl.rclient == nil {
		return ErrUninitialized
	}

	if key == "" {
		return ErrEmptyKey
	}

	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	return rl.tryLock(ctx, key, ttl)
}

// ObtainLock tries to lock a key until try timeout is reached with loopPeriod pause between tries.
func (rl *redisLocker) ObtainLock(
	ctx context.Context,
	key string,
	ttl time.Duration,
	timeout time.Duration,
	retryPeriod time.Duration,
) Error {

	if rl == nil || rl.rclient == nil {
		return ErrUninitialized
	}

	if key == "" {
		return ErrEmptyKey
	}

	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	timeoutTime := now.Add(timeout)
	// use now variable and !now.After condition to make sure at least one try is done
	// even in case 0 timeout is received
	for ; !now.After(timeoutTime); now = time.Now() {
		tryRes := rl.tryLock(ctx, key, ttl)
		switch tryRes {
		case ErrAlreadyLocked:
			time.Sleep(retryPeriod)
		default:
			return tryRes
		}
	}
	// if we reached this point, it means timeout is reached and another lock still not released
	return ErrAlreadyLocked
}

func (rl *redisLocker) Unlock() Error {

	if rl == nil || rl.rclient == nil {
		return ErrUninitialized
	}

	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	if rl.locked != nil && rl.locked.Key != "" {
		rl.rclient.Del(rl.locked.Ctx, rl.locked.Key)
	}
	rl.locked = nil

	return nil
}

// tryLock actually locks the record.
// ! No sync.
// ! No parameters validation.
func (rl *redisLocker) tryLock(ctx context.Context, key string, ttl time.Duration) Error {

	lockKey := lockKeyPrefix + key
	if rl.locked != nil {
		if rl.locked.Key != lockKey {
			return ErrUnlockRequired
		}
		return nil
	}

	if setRes := rl.rclient.SetNX(ctx, lockKey, lockValue, ttl); setRes.Val() {
		rl.locked = &lockInfo{
			Key: lockKey,
			Ctx: ctx,
		}
		return nil
	}

	return ErrAlreadyLocked
}
