package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/felipeascari/kv-store/pkg/lock"
)

type LockedStore struct {
	store              Store
	lockManager        *lock.Manager
	lastProcessedToken map[string]int64
	mu                 sync.RWMutex
}

func NewLockedStore(store Store, lockMgr *lock.Manager) *LockedStore {
	return &LockedStore{
		store:              store,
		lockManager:        lockMgr,
		lastProcessedToken: make(map[string]int64),
	}
}

func (ls *LockedStore) Save(key string, value any) error {
	return ls.lockManager.ExecuteWithLock(context.Background(), key, func(token int64) error {
		if !ls.validateToken(key, token) {
			return fmt.Errorf("token %d rejected: a newer token already processed key %q: %w", token, key, ErrInvalidToken)
		}

		if err := ls.store.Save(key, value); err != nil {
			return fmt.Errorf("failed to save with fencing token %d: %w", token, err)
		}

		ls.recordToken(key, token)
		return nil
	})
}

func (ls *LockedStore) Retrieve(key string) (any, error) {
	var result any

	err := ls.lockManager.ExecuteWithLock(context.Background(), key, func(token int64) error {
		value, err := ls.store.Retrieve(key)
		if err != nil {
			return fmt.Errorf("failed to retrieve with fencing token %d: %w", token, err)
		}

		result = value
		ls.recordToken(key, token)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ls *LockedStore) Delete(key string) error {
	return ls.lockManager.ExecuteWithLock(context.Background(), key, func(token int64) error {
		if !ls.validateToken(key, token) {
			return fmt.Errorf("token %d rejected: a newer token already processed key %q: %w", token, key, ErrInvalidToken)
		}

		if err := ls.store.Delete(key); err != nil {
			return fmt.Errorf("failed to delete with fencing token %d: %w", token, err)
		}

		ls.resetToken(key)
		return nil
	})
}

func (ls *LockedStore) validateToken(key string, token int64) bool {
	ls.mu.RLock()
	defer ls.mu.RUnlock()
	lastToken, exists := ls.lastProcessedToken[key]
	if !exists {
		return true
	}
	return token > lastToken
}

func (ls *LockedStore) recordToken(key string, token int64) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	currentMax, exists := ls.lastProcessedToken[key]
	if !exists || token > currentMax {
		ls.lastProcessedToken[key] = token
	}
}

func (ls *LockedStore) resetToken(key string) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	delete(ls.lastProcessedToken, key)
}
