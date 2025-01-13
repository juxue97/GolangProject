package types

import "time"

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	TTL      time.Duration
	Enabled  bool
}
