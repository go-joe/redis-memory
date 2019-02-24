package redis

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/go-joe/joe"
)

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

func Memory(addr string, opts ...Option) joe.Option {
	return func(joeConf *joe.Config) error {
		conf := Config{Addr: addr}
		for _, opt := range opts {
			err := opt(&conf)
			if err != nil {
				return err
			}
		}

		if conf.Logger == nil {
			conf.Logger = joeConf.Logger
		}

		memory, err := NewMemory(conf)
		if err != nil {
			return err
		}

		joeConf.Memory = memory
		return nil
	}
}

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
		return nil, errors.Wrap(err, "failed to ping redis")
	}

	memory.logger.Info("Memory initialized successfully")
	return memory, nil
}

func (b *memory) Set(key, value string) error {
	resp := b.Client.HSet(b.hkey, key, value)
	return resp.Err()
}

func (b *memory) Get(key string) (string, bool, error) {
	res, err := b.Client.HGet(b.hkey, key).Result()
	switch {
	case err == redis.Nil:
		return "", false, nil
	case err != nil:
		return "", false, err
	default:
		return res, true, nil
	}
}

func (b *memory) Delete(key string) (bool, error) {
	res, err := b.Client.HDel(b.hkey, key).Result()
	return res > 0, err
}

func (b *memory) Memories() (map[string]string, error) {
	return b.Client.HGetAll(b.hkey).Result()
}

func (b *memory) Close() error {
	return b.Client.Close()
}
