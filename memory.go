// Package redis implements a Redis integration for the Joe bot library.
// https://github.com/go-joe/joe
package redis

import (
	"fmt"

	"github.com/go-joe/joe"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

// Config contains all settings for the Redis memory.
type Config struct {
	Addr     string
	Key      string
	Password string
	DB       int
	Logger   *zap.Logger
}

type memory struct {
	logger *zap.Logger
	Client *redis.Client
	hkey   string
}

// Memory returns a jos Module that configures the bot to use Redis as key-value
// store.
func Memory(addr string, opts ...Option) joe.Module {
	return joe.ModuleFunc(func(joeConf *joe.Config) error {
		conf := Config{Addr: addr}
		for _, opt := range opts {
			err := opt(&conf)
			if err != nil {
				return err
			}
		}

		if conf.Logger == nil {
			conf.Logger = joeConf.Logger("redis")
		}

		memory, err := NewMemory(conf)
		if err != nil {
			return err
		}

		joeConf.SetMemory(memory)
		return nil
	})
}

// NewMemory creates a Redis implementation of a joe.Memory.
func NewMemory(conf Config) (joe.Memory, error) {
	if conf.Logger == nil {
		conf.Logger = zap.NewNop()
	}

	if conf.Key == "" {
		conf.Key = "joe-bot"
	}

	memory := &memory{
		logger: conf.Logger,
		hkey:   conf.Key,
	}

	memory.logger.Debug("Connecting to redis memory",
		zap.String("addr", conf.Addr),
		zap.String("key", memory.hkey),
	)

	memory.Client = redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.DB,
	})

	_, err := memory.Client.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	memory.logger.Info("Memory initialized successfully")
	return memory, nil
}

// Set implements joe.Memory by settings the key to the given value in a Redis
// hash set.
func (b *memory) Set(key string, value []byte) error {
	resp := b.Client.HSet(b.hkey, key, value)
	return resp.Err()
}

// Get implements joe.Memory by retrieving a key from a Redis hash set.
func (b *memory) Get(key string) ([]byte, bool, error) {
	res, err := b.Client.HGet(b.hkey, key).Result()
	switch {
	case err == redis.Nil:
		return nil, false, nil
	case err != nil:
		return nil, false, err
	default:
		return []byte(res), true, nil
	}
}

// Delete implements joe.Memory by deleting the given key from the Redis hash set.
func (b *memory) Delete(key string) (bool, error) {
	res, err := b.Client.HDel(b.hkey, key).Result()
	return res > 0, err
}

// Keys implements joe.Memory by returning all previously set keys from Redis
func (b *memory) Keys() ([]string, error) {
	return b.Client.HKeys(b.hkey).Result()
}

// Close terminates the Redis connection
func (b *memory) Close() error {
	return b.Client.Close()
}
