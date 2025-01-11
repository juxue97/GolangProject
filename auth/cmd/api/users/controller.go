package users

import (
	"net/http"

	"github.com/go-chi/chi/v5" // remember to add v5
	"github.com/juxue97/auth/internal/repository"
	"github.com/juxue97/common"
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
//	@Router			/users/activate/{token} [put]
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
