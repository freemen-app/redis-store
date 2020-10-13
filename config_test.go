package redisStore_test

import (
	"errors"
	"testing"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"

	redisStore "github.com/freemen-app/redis-store"
)

func TestConfig_Validate(t *testing.T) {
	type fields struct {
		host        string
		port        string
		database    int
		password    string
		connTimeout time.Duration
		poolSize    uint
	}
	tests := []struct {
		name       string
		fields     fields
		wantErrKey string
	}{
		{
			name: "valid",
			fields: fields{
				host:        "localhost",
				port:        "6379",
				database:    0,
				password:    "test",
				connTimeout: time.Nanosecond,
				poolSize:    1,
			},
		},
		{
			name: "required host",
			fields: fields{
				port:        "6379",
				database:    0,
				password:    "test",
				connTimeout: time.Nanosecond,
				poolSize:    1,
			},
			wantErrKey: "host",
		},
		{
			name: "required port",
			fields: fields{
				host:        "localhost",
				database:    0,
				password:    "test",
				connTimeout: time.Nanosecond,
				poolSize:    1,
			},
			wantErrKey: "port",
		},
		{
			name: "required password",
			fields: fields{
				host:        "localhost",
				port:        "6379",
				database:    0,
				connTimeout: time.Nanosecond,
				poolSize:    1,
			},
			wantErrKey: "password",
		},
		{
			name: "required timeout",
			fields: fields{
				host:     "localhost",
				port:     "6379",
				database: 0,
				password: "test",
				poolSize: 1,
			},
			wantErrKey: "conn_timeout",
		},
		{
			name: "invalid host",
			fields: fields{
				host:        "test@gmail.com",
				port:        "8000",
				database:    0,
				password:    "test",
				connTimeout: time.Nanosecond,
				poolSize:    1,
			},
			wantErrKey: "host",
		},
		{
			name: "invalid port",
			fields: fields{
				host:        "localhost",
				port:        "999999999",
				database:    0,
				password:    "test",
				connTimeout: time.Nanosecond,
				poolSize:    1,
			},
			wantErrKey: "port",
		},
		{
			name: "invalid database",
			fields: fields{
				host:        "localhost",
				port:        "8000",
				database:    -10,
				password:    "test",
				connTimeout: time.Nanosecond,
				poolSize:    1,
			},
			wantErrKey: "database",
		},
		{
			name: "invalid pool size",
			fields: fields{
				host:        "localhost",
				port:        "8000",
				database:    10,
				password:    "test",
				connTimeout: time.Nanosecond,
				poolSize:    0,
			},
			wantErrKey: "pool_size",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := redisStore.Config{
				Host:        tt.fields.host,
				Port:        tt.fields.port,
				Database:    tt.fields.database,
				Password:    tt.fields.password,
				ConnTimeout: tt.fields.connTimeout,
				PoolSize:    tt.fields.poolSize,
			}
			err := c.Validate()
			if tt.wantErrKey == "" {
				assert.Nil(t, err, err)
			} else {
				var validationErr validation.Errors
				assert.True(t, errors.As(err, &validationErr), err)
				assert.Contains(t, validationErr, tt.wantErrKey)
			}
		})
	}
}

func TestConfig_DSN(t *testing.T) {
	conf := redisStore.Config{Host: "localhost", Port: "6379"}
	got := conf.DSN()
	want := "localhost:6379"
	assert.EqualValues(t, got, want)
}
