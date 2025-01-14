package types

import "time"

type RateLimitConfig struct {
	DB      int
	Limit   int
	Window  time.Duration
	Enabled bool
}
