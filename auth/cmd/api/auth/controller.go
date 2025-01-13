package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/juxue97/auth/internal/authenticator"
	"github.com/juxue97/auth/internal/config"
	"github.com/juxue97/auth/internal/mailer"
	"github.com/juxue97/auth/internal/repository"
	"github.com/juxue97/common"
)

type userWithToken struct {
	User  *repository.User
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
		return
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
		return
	}

	// send activation email
	ctx := r.Context()
	plainToken := uuid.New().String()

	// hash the token for storage but keep the plain token for email
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	err := repository.Store.Users.CreateAndInvite(ctx, user, hashToken, config.Configs.Mail.Exp)
	if err != nil {
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
		return
	}

	common.Logger.Infow("comfirmation email sent", "status code", status)

	if err := common.WriteJSON(w, http.StatusCreated, userWithToken); err != nil {
		common.InternalServerError(w, r, err)
		return
	}
}

// createTokenHandler godoc
//
//	@Summary		Login user
//	@Description	Creates a token after successful login
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		loginUserPayload	true	"User credentials"
//	@Success		200		{string}	string					"Token"
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Router			/auth/login [post]
func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload loginUserPayload
	if err := common.ReadJSON(w, r, &payload); err != nil {
		common.BadRequestResponse(w, r, err)
		return
	}
	if err := common.Validate.Struct(payload); err != nil {
		common.BadRequestResponse(w, r, err)
		return
	}
	user, err := repository.Store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch err {
		case common.ErrNotFound:
			common.UnauthorizedError(w, r, err)
		default:
			common.InternalServerError(w, r, err)
		}
		return
	}
	match := user.Password.ComparePassword(payload.Password)
	if !match {
		common.Logger.Errorw("failed to compare password")
		common.UnauthorizedError(w, r, fmt.Errorf("failed to compare password"))
		return
	}

	// Generate JWT token, encode the user ID in it
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(config.Configs.Auth.Token.Exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": config.Configs.Auth.Token.Iss,
		"aud": config.Configs.Auth.Token.Aud,
	}

	token, err := authenticator.JwtAuthenticator.GenerateToken(claims)
	if err != nil {
		common.Logger.Errorw("failed to generate token", "error", err)
		common.InternalServerError(w, r, err)
		return
	}
	if err := common.WriteJSON(w, http.StatusCreated, token); err != nil {
		common.InternalServerError(w, r, err)
	}
}
