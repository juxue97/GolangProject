package test_user

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/juxue97/auth/cmd/api/test_utils"
)

func TestUser(t *testing.T) {
	app := test_utils.NewTestApplication(t)
	mux := app.Mount()

	t.Run("Activate User Tests", func(t *testing.T) {
		tests := []struct {
			name           string
			token          string
			expectedStatus int
		}{
			{
				name:           "Valid activate",
				token:          "valid-token",
				expectedStatus: http.StatusNoContent,
			},
			{
				name:           "Invalid activate",
				token:          "invalid-token",
				expectedStatus: http.StatusNotFound,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				req := test_utils.TestRequestWithoutPayload(
					t,
					http.MethodPut,
					fmt.Sprintf("/v1/users/activate/%v", tc.token),
					tc.expectedStatus,
				)

				rr := test_utils.ExecuteRequest(req, mux)
				test_utils.CheckResponseCode(t, tc.expectedStatus, rr.Code)
			})
		}
	})

	t.Run("JWT Middleware Test", func(t *testing.T) {
		tests := []struct {
			name           string
			id             string
			role           string
			content        string
			url            string
			method         string
			expectedStatus int
		}{
			{
				name:           "Should get unauthorized error",
				id:             "22",
				role:           "admin",
				url:            "/v1/users/",
				method:         http.MethodGet,
				expectedStatus: http.StatusUnauthorized,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				req := test_utils.TestRequestWithoutPayload(
					t,
					http.MethodPost,
					tc.url,
					tc.expectedStatus,
				)

				rr := test_utils.ExecuteRequest(req, mux)
				test_utils.CheckResponseCode(t, tc.expectedStatus, rr.Code)
			})
		}
	})
}
