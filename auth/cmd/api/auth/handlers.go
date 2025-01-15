package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/juxue97/auth/config"
	"github.com/juxue97/auth/internal"
	"github.com/juxue97/auth/internal/authenticator"
	"github.com/juxue97/auth/internal/mailer"
	"github.com/juxue97/auth/internal/repository"
	"github.com/juxue97/common"
	"go.uber.org/zap"
)

type authHandler struct {
	cfg           *config.Config
	logger        *zap.SugaredLogger
	store         *repository.Repository
	authenticator *authenticator.Authenticator
	mailer        *mailer.Client
}

func NewAuthHandler(
	cfg *config.Config,
	logger *zap.SugaredLogger,
	store *repository.Repository,
	authenticator *authenticator.Authenticator,
	mailer *mailer.Client,
) *authHandler {
	return &authHandler{
		cfg:           cfg,
		logger:        logger,
		store:         store,
		authenticator: authenticator,
		mailer:        mailer,
	}
}

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
// @Param			payload	body		RegisterUserPayload	true	"User credentials"
// @Success		201		{object}	userWithToken		"User registered"
// @Failure		400		{object}	error
// @Failure		500		{object}	error
// @Router			/auth/user [post]
func (a *authHandler) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload

	if err := common.ReadJSON(w, r, &payload); err != nil {
		common.BadRequestResponse(w, r, err)
		return
	}
	if err := internal.Validate.Struct(payload); err != nil {
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

	err := a.store.Users.CreateAndInvite(ctx, user, hashToken, a.cfg.Mail.Exp)
	if err != nil {
		switch err {
		case common.ErrUserAlreadyExists:
			common.DuplicateErrorResponse(w, r, err)
		case common.ErrEmailAlreadyExists:
			common.DuplicateErrorResponse(w, r, err)
		default:
			common.InternalServerError(w, r, err)
		}
		return
	}
	userWithToken := userWithToken{
		User:  user,
		Token: plainToken,
	}

	activationURL := fmt.Sprintf("%s/confirm/%s", a.cfg.FrontendURL, plainToken)

	isProdEnv := a.cfg.Env == "production"

	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}
	status, err := a.mailer.MailTrapService.Send(mailer.UserWelcomeInvitationTemplate, user.Username, user.Email, vars, !isProdEnv)
	if err != nil {
		a.logger.Errorw("failed to send welcoming email", "error", err)

		// rollback transaction
		if err := a.store.Users.Delete(ctx, user.ID); err != nil {
			a.logger.Errorw("failed to delete user", "error", err)
		}
		common.InternalServerError(w, r, err)
		return
	}
	a.logger.Infow("comfirmation email sent", "status code", status)

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
//	@Param			payload	body		LoginUserPayload	true	"User credentials"
//	@Success		200		{string}	string					"Token"
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Router			/auth/login [post]
func (a *authHandler) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload LoginUserPayload

	if err := common.ReadJSON(w, r, &payload); err != nil {
		common.BadRequestResponse(w, r, err)
		return
	}

	if err := internal.Validate.Struct(payload); err != nil {
		common.BadRequestResponse(w, r, err)
		return
	}

	user, err := a.store.Users.GetByEmail(r.Context(), payload.Email)
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
		a.logger.Errorw("failed to compare password")
		common.UnauthorizedError(w, r, fmt.Errorf("invalid credentials"))
		return
	}

	// Generate JWT token, encode the user ID in it
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(a.cfg.Auth.Token.Exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": a.cfg.Auth.Token.Iss,
		"aud": a.cfg.Auth.Token.Aud,
	}

	token, err := a.authenticator.Jwt.GenerateToken(claims)
	if err != nil {
		a.logger.Errorw("failed to generate token", "error", err)
		common.InternalServerError(w, r, err)
		return
	}
	if err := common.WriteJSON(w, http.StatusOK, token); err != nil {
		common.InternalServerError(w, r, err)
	}
}
