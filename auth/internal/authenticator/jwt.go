package authenticator

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/juxue97/auth/internal/config"
	"github.com/juxue97/common"
)

type JwtAuth struct {
	secret string
	aud    string
	iss    string
}

func NewJwtAuth(secret string, aud string, iss string) *JwtAuth {
	return &JwtAuth{
		secret: secret,
		iss:    iss,
		aud:    aud,
	}
}

var JwtAuthenticator *JwtAuth

func init() {
	JwtAuthenticator = NewJwtAuth(
		config.Configs.Auth.Token.Secret,
		config.Configs.Auth.Token.Iss,
		config.Configs.Auth.Token.Aud,
	)
	common.Logger.Info("JwtAuthenticator initialized")
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
		jwt.WithValidMethods([]string{jwt.SigningMethodES256.Name}),
	)
}
