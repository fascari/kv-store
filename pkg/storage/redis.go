package storage

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedis(addr, password string, db int) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, err
	}

	return &Redis{
		client: client,
		ctx:    ctx,
	}, nil
}

func (s *Redis) Save(key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.client.Set(s.ctx, key, data, 0).Err()
}

func (s *Redis) Retrieve(key string) (any, error) {
	data, err := s.client.Get(s.ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, ErrKeyNotFound
	}
	if err != nil {
		return nil, err
	}

	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		return nil, err
	}

	return value, nil
}

func (s *Redis) Delete(key string) error {
	result, err := s.client.Del(s.ctx, key).Result()
	if err != nil {
		return err
	}

	if result == 0 {
		return ErrKeyNotFound
	}

	return nil
}

func (s *Redis) Close() error {
	return s.client.Close()
}

func (s *Redis) Ping(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(s.ctx, timeout)
	defer cancel()
	return s.client.Ping(ctx).Err()
}
