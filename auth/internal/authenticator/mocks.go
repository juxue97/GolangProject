package authenticator

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MockJwt struct{}

func NewMockAuthenticator() *Authenticator {
	return &Authenticator{
		Jwt: &MockJwt{},
	}
}

const secret = "test"

var MockClaims = jwt.MapClaims{
	"sub": int64(22),
	"iss": "https://example.com",
	"aud": "https://example.com",
	"exp": time.Now().Add(time.Hour).Unix(),
}

func (a *MockJwt) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MockClaims)

	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString, nil
}

func (a *MockJwt) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
}
