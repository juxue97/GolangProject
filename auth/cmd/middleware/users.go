package middlewares

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/juxue97/auth/internal/repository"
	"github.com/juxue97/common"
)

type (
	userKey       string
	targetUserKey string
)

const (
	userCtx       userKey       = "user"
	targetUserCtx targetUserKey = "targetUser"
)

func (s *MiddlewareService) UsersContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			common.InternalServerError(w, r, err)
			return
		}

		ctx := r.Context()
		user, err := s.getUser(ctx, id)
		if err != nil {
			common.ForbiddenError(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, targetUserCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *MiddlewareService) GetUserFromContext(r *http.Request) (*repository.User, error) {
	user, ok := r.Context().Value(userCtx).(*repository.User)
	if !ok {
		return nil, common.ErrContextNotFound
	}
	return user, nil
}

func (s *MiddlewareService) GetTargetUserFromContext(r *http.Request) (*repository.User, error) {
	user, ok := r.Context().Value(targetUserCtx).(*repository.User)
	if !ok {
		return nil, common.ErrContextNotFound
	}
	return user, nil
}
