package authenticator

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtAuth struct {
	exp    time.Duration
	secret string
	aud    string
	iss    string
}

func NewJwtAuth(exp time.Duration, secret string, aud string, iss string) *Authenticator {
	return &Authenticator{
		Jwt: &JwtAuth{
			exp:    exp,
			secret: secret,
			iss:    iss,
			aud:    aud,
		},
	}
}

func (j *JwtAuth) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *JwtAuth) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(j.secret), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(j.aud),
		jwt.WithIssuer(j.iss),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}
