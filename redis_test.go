package redisStore_test

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	redisStore "redis-store"
)

var (
	conf *redisStore.Config
)

func TestMain(m *testing.M) {
	database, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	conf = &redisStore.Config{
		Host:        os.Getenv("REDIS_HOST"),
		Port:        os.Getenv("REDIS_PORT"),
		Password:    os.Getenv("REDIS_PASSWORD"),
		Database:    database,
		ConnTimeout: time.Second / 2,
		PoolSize:    1,
	}
	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	conf := &redisStore.Config{
		Host:     "localhost",
		Port:     "6379",
		Password: "test",
	}
	store := redisStore.New(conf)
	assert.False(t, store.IsRunning())
}

func TestStore_Start(t *testing.T) {
	type args struct {
		conf *redisStore.Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "succeed",
			args: args{conf: conf},
		},
		{
			name: "invalid credentials",
			args: args{
				conf: &redisStore.Config{
					Host:        conf.Host,
					Port:        conf.Port,
					Password:    "wrong",
					ConnTimeout: time.Second,
					PoolSize:    2,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid host/port",
			args: args{
				conf: &redisStore.Config{
					Host:        "localhost",
					Port:        "5673",
					Password:    "test",
					ConnTimeout: time.Second,
					PoolSize:    2,
				},
			},
			wantErr: true,
		},
		{
			name: "timeout",
			args: args{
				conf: &redisStore.Config{
					Host:        conf.Host,
					Port:        conf.Port,
					Password:    conf.Password,
					ConnTimeout: time.Nanosecond,
					PoolSize:    conf.PoolSize,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := redisStore.New(tt.args.conf)
			gotErr := store.Start()
			assert.EqualValues(t, tt.wantErr, gotErr != nil, gotErr)
			assert.EqualValues(t, store.IsRunning(), gotErr == nil)
			t.Cleanup(store.Shutdown)
		})
	}
}

func Test_store_Shutdown(t *testing.T) {
	store := redisStore.New(conf)

	err := store.Start()
	assert.NoError(t, err)
	assert.True(t, store.IsRunning())

	store.Shutdown()
	assert.False(t, store.IsRunning())
}
