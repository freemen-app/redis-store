package redisStore

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Config struct {
	Host        string        `json:"host"`
	Port        string        `json:"port"`
	Password    string        `json:"password"`
	Database    int           `json:"database"`
	ConnTimeout time.Duration `config:"conn_timeout" json:"conn_timeout"`
	PoolSize    uint          `config:"pool_size" json:"pool_size"`
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.Host, validation.Required, is.Host),
		validation.Field(&c.Port, validation.Required, is.Port),
		validation.Field(&c.Database, validation.Min(0)),
		validation.Field(&c.Password, validation.Required),
		validation.Field(&c.ConnTimeout, validation.Required),
		validation.Field(&c.PoolSize, validation.Required),
	)
}

func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
