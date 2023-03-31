package redislock

import (
	"errors"
)

type RedisLockError error

var (
	// these errors describes a result of a record lock
	Ok             = RedisLockError(nil)
	AlreadyLocked  = RedisLockError(errors.New("key is already locked"))
	Uninitialized  = RedisLockError(errors.New("locker or redis client is not initialized"))
	EmptyKey       = RedisLockError(errors.New("key to lock is empty"))
	UnlockRequired = RedisLockError(errors.New("locker already locked another key"))
)
