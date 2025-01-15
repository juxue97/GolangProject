package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/juxue97/auth/internal/repository"
	"github.com/juxue97/common"
)

// JWT Middleware to validate token
func (s *MiddlewareService) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the authorization header
		getHeader := r.Header.Get("Authorization")
		if getHeader == "" {
			common.UnauthorizedMiddlewareError(w, r, fmt.Errorf("missing authorization header"))
			return
		}

		parts := strings.Split(getHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			common.UnauthorizedMiddlewareError(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		// Get the token
		token := parts[1]
		jwtToken, err := s.authenticator.Jwt.ValidateToken(token)
		if err != nil {
			common.UnauthorizedMiddlewareError(w, r, err)
			return
		}

		// Get the user
		claims, _ := jwtToken.Claims.(jwt.MapClaims)
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			common.UnauthorizedMiddlewareError(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := s.getUser(ctx, userID)
		if err != nil {
			common.UnauthorizedMiddlewareError(w, r, err)
			return
		}
		// Add the user to the context
		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *MiddlewareService) getUser(ctx context.Context, userID int64) (*repository.User, error) {
	// if redis is not enabled, retrieve from database
	if !s.cfg.RedisCfg.Enabled {
		return s.PgStore.Users.GetByID(ctx, userID)
	}

	user, err := s.cacheStorage.Users.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	// if no records is found, retrieve from database, then set to redis
	var userDB *repository.User
	if user == nil {

		userDB, err = s.PgStore.Users.GetByID(ctx, userID)
		if err != nil {
			return nil, err
		}
		if err := s.cacheStorage.Users.Set(ctx, userDB); err != nil {
		}
		return userDB, err

	}
	return user, nil
}
