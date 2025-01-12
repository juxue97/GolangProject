package authenticator

import "github.com/golang-jwt/jwt/v5"

type Jwt interface {
	GenerateToken(jwt.Claims) (string, error)
	ValidateToken(string) (*jwt.Token, error)
}
