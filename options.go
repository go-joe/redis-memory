package redis

import "go.uber.org/zap"

// An Option can be passed to the Memory function for opt-in functionality.
type Option func(*Config) error

// WithConfig is an Option to have full control over all redis connection options.
func WithConfig(newConf Config) Option {
	return func(oldConf *Config) error {
		oldConf.Addr = newConf.Addr
		oldConf.Key = newConf.Key
		oldConf.Password = newConf.Password
		oldConf.DB = newConf.DB
		oldConf.Logger = newConf.Logger
		return nil
	}
}

// WithLogger is an Option to let the Redis memory use a specific logger.
func WithLogger(logger *zap.Logger) Option {
	return func(conf *Config) error {
		conf.Logger = logger
		return nil
	}
}

// WithKey is an Option to use a different redis key to store the memories of
// the bot (default is "joe-bot").
func WithKey(key string) Option {
	return func(conf *Config) error {
		conf.Key = key
		return nil
	}
}
