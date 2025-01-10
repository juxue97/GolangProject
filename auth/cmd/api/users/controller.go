package users

import (
	"fmt"
	"net/http"

	"github.com/juxue97/common"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserPayload
	if err := common.ReadJSON(w, r, &payload); err != nil {
		common.BadRequestResponse(w, r, err)
	}
	fmt.Println(payload)
}
