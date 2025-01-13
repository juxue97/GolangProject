package types

import "time"

type RateLimitConfig struct {
	Limit   int
	Window  time.Duration
	Enabled bool
}
