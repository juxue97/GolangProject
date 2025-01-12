package types

import "time"

type AuthConfig struct {
	Token TokenConfig
}

type TokenConfig struct {
	Secret string
	Exp    time.Duration
	Iss    string
	Aud    string
}
