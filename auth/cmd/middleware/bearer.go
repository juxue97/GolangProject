package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/juxue97/auth/cmd/api/users"
	"github.com/juxue97/auth/internal/authenticator"
	"github.com/juxue97/auth/internal/cache"
	"github.com/juxue97/auth/internal/config"
	"github.com/juxue97/auth/internal/repository"
	"github.com/juxue97/common"
)

// JWT Middleware to validate token
func AuthTokenMiddleware(next http.Handler) http.Handler {
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
		jwtToken, err := authenticator.JwtAuthenticator.ValidateToken(token)
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

		user, err := getUser(ctx, userID)
		if err != nil {
			common.UnauthorizedMiddlewareError(w, r, err)
			return
		}
		// Add the user to the context
		ctx = context.WithValue(ctx, users.UserCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUser(ctx context.Context, userID int64) (*repository.User, error) {
	// if redis is not enabled, retrieve from database
	if !config.Configs.RedisCfg.Enabled {
		return repository.Store.Users.GetByID(ctx, userID)
	}

	user, err := cache.CacheStorage.Users.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	// if no records is found, retrieve from database, then set to redis
	if user == nil {
		user, err := repository.Store.Users.GetByID(ctx, userID)
		if err != nil {
			return nil, err
		}
		if err := cache.CacheStorage.Users.Set(ctx, user); err != nil {
			return nil, err
		}
	}

	return user, nil
}
