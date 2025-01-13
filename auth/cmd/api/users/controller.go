package users

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5" // remember to add v5
	"github.com/juxue97/auth/internal/cache"
	"github.com/juxue97/auth/internal/repository"
	"github.com/juxue97/common"
)

type (
	userKey       string
	targetUserKey string
)

const (
	UserCtx       userKey       = "user"
	TargetUserCtx targetUserKey = "targetUser"
)

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
func ActivateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	err := repository.Store.Users.ActivateUser(r.Context(), token)
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
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := repository.Store.Users.GetAll(r.Context())
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
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// get user from context, already authenticated and found the target id user
	userCtx, err := GetTargetUserFromContext(r)
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
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload UpdateUserPayload
	if err := common.ReadJSON(w, r, &payload); err != nil {
		common.BadRequestResponse(w, r, err)
		return
	}

	if err := common.Validate.Struct(payload); err != nil {
		common.BadRequestResponse(w, r, err)
		return
	}

	userCtx, err := GetTargetUserFromContext(r)
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
	if err := repository.Store.Users.Update(ctx, user); err != nil {
		switch err {
		case sql.ErrNoRows:
			common.NotFoundError(w, r, err)
		default:
			common.InternalServerError(w, r, err)
		}
		return
	}

	// If success update, clear the cache
	err = cache.CacheStorage.Users.Delete(ctx, user.ID)
	if err != nil {
		common.InternalServerError(w, r, err)
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
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userCtx, err := GetTargetUserFromContext(r)
	if err != nil {
		common.InternalServerError(w, r, err)
		return
	}

	ctx := r.Context()
	if err := repository.Store.Users.DeleteUser(ctx, userCtx.ID); err != nil {
		switch err {
		case sql.ErrNoRows:
			common.NotFoundError(w, r, err)
		default:
			common.InternalServerError(w, r, err)
		}
		return
	}

	// If success delete, clear the cache
	err = cache.CacheStorage.Users.Delete(ctx, userCtx.ID)
	if err != nil {
		common.InternalServerError(w, r, err)
	}

	if err := common.WriteJSON(w, http.StatusNoContent, nil); err != nil {
		common.InternalServerError(w, r, err)
	}
}

func GetUserFromContext(r *http.Request) (*repository.User, error) {
	user, ok := r.Context().Value(UserCtx).(*repository.User)
	if !ok {
		return nil, common.ErrContextNotFound
	}
	return user, nil
}

func GetTargetUserFromContext(r *http.Request) (*repository.User, error) {
	user, ok := r.Context().Value(TargetUserCtx).(*repository.User)
	if !ok {
		return nil, common.ErrContextNotFound
	}
	return user, nil
}
