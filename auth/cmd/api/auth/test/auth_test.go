package test_auth

import (
	"net/http"
	"testing"

	"github.com/juxue97/auth/cmd/api/auth"
	"github.com/juxue97/auth/cmd/api/test_utils"
)

func TestAuth(t *testing.T) {
	app := test_utils.NewTestApplication(t)
	mux := app.Mount()

	t.Run("Login Tests", func(t *testing.T) {
		tests := []struct {
			name           string
			payload        auth.LoginUserPayload
			expectedStatus int
		}{
			{
				name: "Valid Login",
				payload: auth.LoginUserPayload{
					Email:    "hwteh1997@hotmail.com",
					Password: "123456",
				},
				expectedStatus: http.StatusOK,
			},
			{
				name: "Invalid Login",
				payload: auth.LoginUserPayload{
					Email:    "hwteh1997@hotmail.com",
					Password: "1234567",
				},
				expectedStatus: http.StatusUnauthorized,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				req := test_utils.TestRequestWithPayload(
					t,
					http.MethodPost,
					"/v1/auth/login",
					tc.payload,
					tc.expectedStatus,
				)

				rr := test_utils.ExecuteRequest(req, mux)
				test_utils.CheckResponseCode(t, tc.expectedStatus, rr.Code)
			})
		}
	})

	t.Run("Register Tests", func(t *testing.T) {
		tests := []struct {
			name           string
			payload        auth.RegisterUserPayload
			expectedStatus int
		}{
			{
				name: "Valid Registration",
				payload: auth.RegisterUserPayload{
					Username: "imGoot",
					Email:    "goot@gmail.com",
					Password: "veryGootPass",
				},
				expectedStatus: http.StatusCreated,
			},
			{
				name: "Invalid Login-Duplicate email",
				payload: auth.RegisterUserPayload{
					Username: "imGoot",
					Email:    "hwteh1997@gmail.com",
					Password: "veryGootPass",
				},
				expectedStatus: http.StatusConflict,
			},
			{
				name: "Invalid Login-Duplicate username",
				payload: auth.RegisterUserPayload{
					Username: "MehNohNah",
					Email:    "tehwei123@gmail.com",
					Password: "veryGootPass",
				},
				expectedStatus: http.StatusConflict,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				req := test_utils.TestRequestWithPayload(
					t,
					http.MethodPost,
					"/v1/auth/user",
					tc.payload,
					tc.expectedStatus,
				)

				rr := test_utils.ExecuteRequest(req, mux)
				test_utils.CheckResponseCode(t, tc.expectedStatus, rr.Code)
			})
		}
	})
}
