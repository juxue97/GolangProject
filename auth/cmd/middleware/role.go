package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/juxue97/auth/internal/repository"
	"github.com/juxue97/common"
)

// type RoleMiddlewareService struct{}

// func NewRoleMiddlewareService() *RoleMiddlewareService {
// 	return &RoleMiddlewareService{}
// }

// Role Middleware: Check the user role,
func (s *MiddlewareService) RoleMiddleware(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := s.GetUserFromContext(r)
		if err != nil {
			common.ForbiddenError(w, r, err)
			return
		}

		// admin can do anything
		if user.Role.Level >= 3 {
			next.ServeHTTP(w, r)
			return
		}

		targetUser, err := s.GetTargetUserFromContext(r)
		if err != nil {
			common.ForbiddenError(w, r, err)
			return
		}

		// if same user, no need to check
		if user.ID == targetUser.ID {
			next.ServeHTTP(w, r)
			return
		}

		allowed, err := s.checkRolePrecedence(r.Context(), user, requiredRole)
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

func (s *MiddlewareService) checkRolePrecedence(ctx context.Context, user *repository.User, requiredRole string) (bool, error) {
	// Retrieve the required role if from database
	role, err := s.PgStore.Roles.GetByName(ctx, requiredRole)
	if err != nil {
		return false, err
	}

	// Check if the u ser has the required role
	return user.Role.Level >= role.Level, nil
}
