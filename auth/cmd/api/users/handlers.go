package users

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5" // remember to add v5
	middlewares "github.com/juxue97/auth/cmd/middleware"
	"github.com/juxue97/auth/internal"
	"github.com/juxue97/auth/internal/cache"
	"github.com/juxue97/auth/internal/repository"
	"github.com/juxue97/common"
)

type userHandler struct {
	middlewareService *middlewares.MiddlewareService
	store             *repository.Repository
	cacheStorage      *cache.RedisCacheStorage
}

func NewUserHandler(middlewareService *middlewares.MiddlewareService, store *repository.Repository, cacheStorage *cache.RedisCacheStorage) *userHandler {
	return &userHandler{
		middlewareService: middlewareService,
		store:             store,
		cacheStorage:      cacheStorage,
	}
}

// ActivateUser godoc
//
//	@Summary		Activates a user account status
//	@Description	Activates a user by invitation token
//	@Tags			users
//	@Produce		json
//	@Param			token	path		string	true	"Invitation token"
//	@Success		204		{string}	string	"User activated"
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//
// @Router			/users/activate/{token} [put]
func (u *userHandler) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	err := u.store.Users.ActivateUser(r.Context(), token)
	if err != nil {
		switch err {
		case common.ErrNotFound:
			common.NotFoundError(w, r, err)
		default:
			common.InternalServerError(w, r, err)
		}
		return
	}
	if err := common.WriteJSON(w, http.StatusNoContent, ""); err != nil {
		common.InternalServerError(w, r, err)
	}
}

// GetAllUsers godoc
//
//	@Summary		Retrieve all users
//	@Description	Retrieve all users
//	@Tags			users
//	@Produce		json
//
// @Success        200     {array}     repository.User
//
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//
// @Security ApiKeyAuth
//
//	@Router			/users/ [get]
func (u *userHandler) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := u.store.Users.GetAll(r.Context())
	if err != nil {
		switch err {
		case common.ErrNotFound:
			common.NotFoundError(w, r, err)
		default:
			common.InternalServerError(w, r, err)
		}
		return
	}
	if err := common.WriteJSON(w, http.StatusOK, users); err != nil {
		common.InternalServerError(w, r, err)
	}
}

// GetAUser godoc
//
//	@Summary		Retrieve single user information
//	@Description	Retrieve single user information by id
//	@Tags			users
//	@Produce		json
//	@Param			id	path		string	true	"User id"
//	@Success        200     {object}     repository.User
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//
// @Security ApiKeyAuth
//
//	@Router			/users/{id} [get]
func (u *userHandler) getUserHandler(w http.ResponseWriter, r *http.Request) {
	// get user from context, already authenticated and found the target id user
	userCtx, err := u.middlewareService.GetTargetUserFromContext(r)
	if err != nil {
		common.InternalServerError(w, r, err)
		return
	}

	if err := common.WriteJSON(w, http.StatusOK, userCtx); err != nil {
		common.InternalServerError(w, r, err)
	}
}

// UpdateUser godoc
//
//	@Summary		Update user
//	@Description	Update user by id
//	@Tags			users
//	@Produce		json
//	@Param			id	path		string	true	"User id"
//	@Param			payload	body		UpdateUserPayload	true	"Post payload"
//	@Success		204		{string}	string	"User updated"
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//
// @Security ApiKeyAuth
//
//	@Router			/users/{id} [put]
func (u *userHandler) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload UpdateUserPayload
	if err := common.ReadJSON(w, r, &payload); err != nil {
		common.BadRequestResponse(w, r, err)
		return
	}

	if err := internal.Validate.Struct(payload); err != nil {
		common.BadRequestResponse(w, r, err)
		return
	}

	userCtx, err := u.middlewareService.GetTargetUserFromContext(r)
	if err != nil {
		common.InternalServerError(w, r, err)
		return
	}

	user := &repository.User{
		ID:       userCtx.ID,
		Username: payload.Username,
		Email:    payload.Email,
	}
	ctx := r.Context()
	if err := u.store.Users.Update(ctx, user); err != nil {
		switch err {
		case sql.ErrNoRows:
			common.NotFoundError(w, r, err)
		default:
			common.InternalServerError(w, r, err)
		}
		return
	}

	// If success update, clear the cache
	err = u.cacheStorage.Users.Delete(ctx, user.ID)
	if err != nil {
		common.InternalServerError(w, r, err)
		return
	}

	if err := common.WriteJSON(w, http.StatusNoContent, nil); err != nil {
		common.InternalServerError(w, r, err)
	}
}

// DeleteUser godoc
//
//	@Summary		Delete user
//	@Description	Delete user by id
//	@Tags			users
//	@Produce		json
//	@Param			id	path		string	true	"User id"
//	@Success		200		{string}	string	"User deleted"
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//
// @Security ApiKeyAuth
//
//	@Router			/users/{id} [delete]
func (u *userHandler) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userCtx, err := u.middlewareService.GetTargetUserFromContext(r)
	if err != nil {
		common.InternalServerError(w, r, err)
		return
	}

	ctx := r.Context()
	if err := u.store.Users.DeleteUser(ctx, userCtx.ID); err != nil {
		switch err {
		case sql.ErrNoRows:
			common.NotFoundError(w, r, err)
		default:
			common.InternalServerError(w, r, err)
		}
		return
	}

	// If success delete, clear the cache
	err = u.cacheStorage.Users.Delete(ctx, userCtx.ID)
	if err != nil {
		common.InternalServerError(w, r, err)
	}

	if err := common.WriteJSON(w, http.StatusNoContent, nil); err != nil {
		common.InternalServerError(w, r, err)
	}
}
