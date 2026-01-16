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

func (r *Redis) Save(key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(r.ctx, key, data, 0).Err()
}

func (r *Redis) Retrieve(key string) (any, error) {
	data, err := r.client.Get(r.ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrKeyNotFound
		}
		return nil, err
	}

	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		return nil, err
	}

	return value, nil
}

func (r *Redis) Delete(key string) error {
	result, err := r.client.Del(r.ctx, key).Result()
	if err != nil {
		return err
	}

	if result == 0 {
		return ErrKeyNotFound
	}

	return nil
}

func (r *Redis) Close() error {
	return r.client.Close()
}

func (r *Redis) Client() *redis.Client {
	return r.client
}

func (r *Redis) Ping(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(r.ctx, timeout)
	defer cancel()
	return r.client.Ping(ctx).Err()
}
