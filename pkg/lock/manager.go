package lock

import (
	"context"
	"fmt"
	"os"
	"time"
)

type (
	Lock interface {
		Acquire(ctx context.Context, key string) (int64, error)
		Release(ctx context.Context, key string, token int64) error
		ValidateToken(ctx context.Context, key string, token int64) (bool, error)
	}

	Manager struct {
		lock           Lock
		serverID       string
		lockTTL        time.Duration
		acquireTimeout time.Duration
	}

	Config struct {
		LockTTL        time.Duration
		AcquireTimeout time.Duration
	}
)

func NewManager(lock Lock) *Manager {
	hostname, _ := os.Hostname()
	return &Manager{
		lock:           lock,
		serverID:       fmt.Sprintf("%s-%d", hostname, os.Getpid()),
		lockTTL:        5 * time.Second,
		acquireTimeout: 30 * time.Second,
	}
}

func (lm *Manager) WithConfig(cfg Config) *Manager {
	if cfg.LockTTL > 0 {
		lm.lockTTL = cfg.LockTTL
	}
	if cfg.AcquireTimeout > 0 {
		lm.acquireTimeout = cfg.AcquireTimeout
	}
	return lm
}

func (lm *Manager) ExecuteWithLock(ctx context.Context, key string, fn func(int64) error) error {
	acquireCtx, cancel := context.WithTimeout(ctx, lm.acquireTimeout)
	defer cancel()

	token, err := lm.lock.Acquire(acquireCtx, key)
	if err != nil {
		return err
	}

	defer func() {
		releaseCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = lm.lock.Release(releaseCtx, key, token)
	}()

	return fn(token)
}
