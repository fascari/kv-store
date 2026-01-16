package lock

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const releaseLockScript = `
	if redis.call('get', KEYS[1]) == ARGV[1] then
		return redis.call('del', KEYS[1])
	else
		return 0
	end
`

var (
	ErrLockAcquisition = errors.New("lock acquisition failed")
	ErrInvalidToken    = errors.New("invalid or expired fencing token")
)

type (
	RedisLock struct {
		client        *redis.Client
		lockKeyPrefix string
		tokenKey      string
		ttl           time.Duration
	}

	Entry struct {
		Token      int64     `json:"token"`
		ServerID   string    `json:"server_id"`
		AcquiredAt time.Time `json:"acquired_at"`
	}
)

func NewRedisLock(client *redis.Client, ttl time.Duration) *RedisLock {
	if ttl == 0 {
		ttl = 5 * time.Second
	}

	return &RedisLock{
		client:        client,
		lockKeyPrefix: "lock:",
		tokenKey:      "lock:token_counter",
		ttl:           ttl,
	}
}

func (rl *RedisLock) Acquire(ctx context.Context, key string) (int64, error) {
	lockKey := rl.lockKeyPrefix + key

	token, err := rl.client.Incr(ctx, rl.tokenKey).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to generate token: %w", err)
	}

	entry := Entry{
		Token:      token,
		ServerID:   "",
		AcquiredAt: time.Now(),
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal lock entry: %w", err)
	}

	result, err := rl.client.SetNX(ctx, lockKey, data, rl.ttl).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to acquire lock: %w", err)
	}

	if !result {
		return 0, fmt.Errorf("lock already held: %w", ErrLockAcquisition)
	}

	return token, nil
}

// Release releases a lock only if the provided token matches.
//
// Uses a Lua script to ensure atomicity: the Get + Compare + Del operations
// must execute as a single atomic operation to prevent race conditions.
// Without atomicity, another process could acquire the lock between the
// validation and deletion, causing this process to delete someone else's lock.
func (rl *RedisLock) Release(ctx context.Context, key string, token int64) error {
	lockKey := rl.lockKeyPrefix + key

	data, err := rl.client.Get(ctx, lockKey).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return fmt.Errorf("failed to get lock: %w", err)
	}

	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return fmt.Errorf("failed to unmarshal lock entry: %w", err)
	}

	if entry.Token != token {
		return fmt.Errorf("token mismatch: expected %d, got %d: %w", token, entry.Token, ErrInvalidToken)
	}

	script := redis.NewScript(releaseLockScript)

	result, err := script.Run(ctx, rl.client, []string{lockKey}, data).Result()
	if err != nil {
		return fmt.Errorf("failed to release lock: %w", err)
	}

	if result.(int64) == 0 {
		return fmt.Errorf("lock token no longer valid: %w", ErrInvalidToken)
	}

	return nil
}

func (rl *RedisLock) ValidateToken(ctx context.Context, key string, token int64) (bool, error) {
	lockKey := rl.lockKeyPrefix + key

	data, err := rl.client.Get(ctx, lockKey).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, fmt.Errorf("failed to get lock: %w", err)
	}

	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return false, fmt.Errorf("failed to unmarshal lock entry: %w", err)
	}

	return entry.Token == token, nil
}
