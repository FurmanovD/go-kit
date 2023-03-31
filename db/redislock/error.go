package redislock

import (
	"errors"
)

type Error error

var (
	// these errors describes a result of a record lock
	ErrAlreadyLocked  = Error(errors.New("key is already locked"))
	ErrUninitialized  = Error(errors.New("locker or redis client is not initialized"))
	ErrEmptyKey       = Error(errors.New("key to lock is empty"))
	ErrUnlockRequired = Error(errors.New("locker already locked another key"))
)
