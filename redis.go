package redisStore

import (
	"github.com/mediocregopher/radix/v3"
	"github.com/rs/zerolog/log"
)

type (
	RedisClient interface {
		Do(action radix.Action) error
	}

	Store interface {
		RedisClient
		Start() error
		IsRunning() bool
		Shutdown()
	}

	store struct {
		config *Config
		*radix.Pool
		isRunning bool
	}
)

func New(config *Config) Store {
	return &store{
		config: config,
	}
}

func (s *store) Start() error {
	customConnFunc := func(network, addr string) (radix.Conn, error) {
		return radix.Dial(
			network, addr,
			radix.DialTimeout(s.config.ConnTimeout),
			radix.DialAuthPass(s.config.Password),
			radix.DialSelectDB(s.config.Database),
		)
	}
	pool, err := radix.NewPool(
		"tcp",
		s.config.DSN(),
		int(s.config.PoolSize),
		radix.PoolConnFunc(customConnFunc),
	)
	if err != nil {
		return err
	}

	s.isRunning = true
	s.Pool = pool
	return nil
}

func (s *store) IsRunning() bool {
	return s.isRunning
}

func (s *store) Shutdown() {
	if s.Pool != nil {
		if err := s.Close(); err != nil {
			log.Error().Msg(err.Error())
		}
	}
	s.isRunning = false
}
