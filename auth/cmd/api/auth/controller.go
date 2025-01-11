package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/juxue97/auth/internal/config"
	"github.com/juxue97/auth/internal/mailer"
	"github.com/juxue97/auth/internal/repository"
	"github.com/juxue97/common"
)

type userWithToken struct {
	*repository.User
	Token string `json:"token"`
}

// registerUserHandler godoc

// @Summary		Registers a user
// @Description	Registers a user and send them an comfirmation email
// @Tags			authentication
// @Accept			json
// @Produce		json
// @Param			payload	body		registerUserPayload	true	"User credentials"
// @Success		201		{object}	userWithToken		"User registered"
// @Failure		400		{object}	error
// @Failure		500		{object}	error
// @Router			/auth/user [post]
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload registerUserPayload
	if err := common.ReadJSON(w, r, &payload); err != nil {
		common.BadRequestResponse(w, r, err)
		return
	}
	if err := common.Validate.Struct(payload); err != nil {
		common.BadRequestResponse(w, r, err)
	}

	user := &repository.User{
		Username: payload.Username,
		Email:    payload.Email,
		Role: repository.Role{
			Name: "user",
		},
	}
	// hash password
	if err := user.Password.SetPassword(payload.Password); err != nil {
		common.InternalServerError(w, r, err)
	}

	// send activation email
	ctx := r.Context()
	plainToken := uuid.New().String()

	// hash the token for storage but keep the plain token for email
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	if err := repository.Store.Users.CreateAndInvite(ctx, user, hashToken, config.Configs.Mail.Exp); err != nil {
		switch err {
		case common.ErrUserAlreadyExists:
			common.BadRequestResponse(w, r, err)
		case common.ErrEmailAlreadyExists:
			common.BadRequestResponse(w, r, err)
		default:
			common.InternalServerError(w, r, err)
		}
		return
	}

	userWithToken := userWithToken{
		User:  user,
		Token: plainToken,
	}

	activationURL := fmt.Sprintf("%s/confirm/%s", config.Configs.FrontendURL, plainToken)

	isProdEnv := config.Configs.Env == "production"

	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}

	status, err := mailer.MailTrapMailer.Send(mailer.UserWelcomeInvitationTemplate, user.Username, user.Email, vars, !isProdEnv)
	if err != nil {
		common.Logger.Errorw("failed to send welcoming email", "error", err)

		// rollback transaction
		if err := repository.Store.Users.Delete(ctx, user.ID); err != nil {
			common.Logger.Errorw("failed to delete user", "error", err)
		}
		common.InternalServerError(w, r, err)
	}

	common.Logger.Infow("comfirmation email sent", "status code", status)

	if err := common.WriteJSON(w, http.StatusCreated, userWithToken); err != nil {
		common.InternalServerError(w, r, err)
	}
}

// func CreateTokenHandler(w http.ResponseWriter, r *http.Request) {
// 	var paylod createTokenPayload
// 	if err := common.ReadJSON(w, r, &paylod); err != nil {
// 		common.BadRequestResponse(w, r, err)
// 	}
// }
