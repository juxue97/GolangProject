package middlewares

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/juxue97/auth/cmd/api/users"
	"github.com/juxue97/common"
)

func UsersContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			common.InternalServerError(w, r, err)
			return
		}

		ctx := r.Context()
		user, err := getUser(ctx, id)
		if err != nil {
			common.ForbiddenError(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, users.TargetUserCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
