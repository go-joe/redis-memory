package redis

import (
	"context"
	"testing"

	"github.com/fzerorubigd/redimock"
	"github.com/go-joe/joe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func newTestMemory(t *testing.T) (joe.Memory, *redimock.Server, func()) {
	ctx, cancel := context.WithCancel(context.Background())

	mock, err := redimock.NewServer(ctx, "")
	require.NoError(t, err)

	conf := Config{
		Addr:   mock.Addr().String(),
		Logger: zaptest.NewLogger(t),
		Key:    "test",
	}

	mock.ExpectPing().Once()
	m, err := NewMemory(conf)
	if !assert.NoError(t, err) {
		cancel()
		t.FailNow()
	}

	return m, mock, cancel
}

func TestMemory(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mock, err := redimock.NewServer(ctx, "")
	require.NoError(t, err)
	mock.ExpectPing().Once()

	logger := zaptest.NewLogger(t)
	addr := mock.Addr().String()
	mod := Memory(addr, WithKey("test-123"), WithLogger(logger))

	store := joe.NewStorage(logger)
	joeConf := joe.NewConfig(logger, nil, store, nil)
	err = mod.Apply(&joeConf)
	require.NoError(t, err)
}

func TestNewMemory_NoRedis(t *testing.T) {
	conf := Config{Addr: ":1"}
	m, err := NewMemory(conf)
	assert.Error(t, err)
	assert.Nil(t, m)
}

func TestNewMemory(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mock, err := redimock.NewServer(ctx, "")
	require.NoError(t, err)
	mock.ExpectPing().Once()

	conf := Config{Addr: mock.Addr().String()}
	m, err := NewMemory(conf)
	require.NoError(t, err)

	mem := m.(*memory)
	assert.Equal(t, "joe-bot", mem.hkey)
	assert.NotNil(t, mem.logger)
}

func TestMemory_Set(t *testing.T) {
	m, mock, cancel := newTestMemory(t)
	defer cancel()

	mock.Expect("HSET").WithArgs("test", "foo", "bar").WillReturn(0).Once()
	assert.NoError(t, m.Set("foo", []byte("bar")))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMemory_Get(t *testing.T) {
	m, mock, cancel := newTestMemory(t)
	defer cancel()

	mock.Expect("HGET").WithArgs("test", "unknown").WillReturn(nil).Once()
	val, ok, err := m.Get("unknown")
	assert.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, val)

	mock.Expect("HGET").WithArgs("test", "foo").WillReturn("bar").Once()
	val, ok, err = m.Get("foo")
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, []byte("bar"), val)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMemory_Get_Error(t *testing.T) {
	m, mock, cancel := newTestMemory(t)
	defer cancel()

	mock.Expect("HGET").WithArgs("test", "error").WillReturn(redimock.Error("this did not work")).Once()
	val, ok, err := m.Get("error")
	assert.EqualError(t, err, "this did not work")
	assert.False(t, ok)
	assert.Nil(t, val)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMemory_Delete(t *testing.T) {
	m, mock, cancel := newTestMemory(t)
	defer cancel()

	mock.Expect("HDEL").WithArgs("test", "unknown").WillReturn(0).Once()
	ok, err := m.Delete("unknown")
	assert.NoError(t, err)
	assert.False(t, ok)

	mock.Expect("HDEL").WithArgs("test", "foo").WillReturn(1).Once()
	ok, err = m.Delete("foo")
	assert.NoError(t, err)
	assert.True(t, ok)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMemory_Keys(t *testing.T) {
	m, mock, cancel := newTestMemory(t)
	defer cancel()

	mock.Expect("HKEYS").WithArgs("test").WillReturn([]string{"foo", "bar"}).Once()
	keys, err := m.Keys()
	assert.NoError(t, err)
	assert.Equal(t, []string{"foo", "bar"}, keys)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMemory_Close(t *testing.T) {
	m, mock, cancel := newTestMemory(t)
	defer cancel()

	mock.ExpectQuit()
	assert.NoError(t, m.Close())
	assert.NoError(t, mock.ExpectationsWereMet())
}
