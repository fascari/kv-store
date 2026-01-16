package storage

import "errors"

const (
	TypeMemory Type = "memory"
	TypeRedis  Type = "redis"
)

var (
	ErrKeyNotFound  = errors.New("key not found")
	ErrInvalidToken = errors.New("invalid fencing token")
)

type (
	Type string

	LockInfo struct {
		Token     int64
		Key       string
		ExpiresAt int64
		ServerID  string
	}
)

func (t Type) String() string {
	return string(t)
}

func (t Type) IsValid() bool {
	switch t {
	case TypeMemory, TypeRedis:
		return true
	default:
		return false
	}
}
