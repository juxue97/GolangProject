package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/juxue97/auth/cmd/api/users"
	"github.com/juxue97/auth/internal/repository"
	"github.com/juxue97/common"
)

// Role Middleware: Check the user role,
func RoleMiddleware(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := users.GetUserFromContext(r)
		targetUser := users.GetTargetUserFromContext(r)

		// if same, no need to check
		if user.ID == targetUser.ID {
			next.ServeHTTP(w, r)
			return
		}

		allowed, err := checkRolePrecedence(r.Context(), user, requiredRole)
		if err != nil {
			common.InternalServerError(w, r, err)
			return
		}
		if !allowed {
			common.ForbiddenError(w, r, fmt.Errorf("forbidden action"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func checkRolePrecedence(ctx context.Context, user *repository.User, requiredRole string) (bool, error) {
	// Retrieve the required role if from database
	role, err := repository.Store.Roles.GetByName(ctx, requiredRole)
	if err != nil {
		return false, err
	}

	// Check if the u ser has the required role
	return user.Role.Level >= role.Level, nil
}
