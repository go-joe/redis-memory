package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestWithConfig(t *testing.T) {
	logger := zaptest.NewLogger(t)
	opt := WithConfig(Config{
		Key:      "new",
		Addr:     ":5678",
		Password: "secret",
		DB:       42,
		Logger:   logger,
	})

	conf := Config{Key: "old", Addr: ":1234"}
	err := opt(&conf)
	require.NoError(t, err)
	assert.Equal(t, "new", conf.Key)
	assert.Equal(t, ":5678", conf.Addr)
	assert.Equal(t, "secret", conf.Password)
	assert.Equal(t, 42, conf.DB)
	assert.Equal(t, logger, conf.Logger)
}
