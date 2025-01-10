package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/juxue97/auth/internal/repository"
	"github.com/juxue97/common"
)

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
	fmt.Println(user.Username, user.Email)
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
	fmt.Println(ctx, hashToken)
	// if err := repository.Store.Users.CreateAndInvite(ctx, user, hashToken, config.Configs.mail.exp); err != nil {
	// }
	common.Logger.Infow("user registered", "user", nil)

	if err := common.WriteJSON(w, http.StatusCreated, nil); err != nil {
		common.InternalServerError(w, r, err)
	}
}

// func CreateTokenHandler(w http.ResponseWriter, r *http.Request) {
// 	var paylod createTokenPayload
// 	if err := common.ReadJSON(w, r, &paylod); err != nil {
// 		common.BadRequestResponse(w, r, err)
// 	}
// }
