package test_auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/juxue97/auth/cmd/api/auth"
	"github.com/juxue97/auth/cmd/api/test_utils"
)

func TestAuth(t *testing.T) {
	app := test_utils.NewTestApplication(t)
	mux := app.Mount()

	t.Run("should be able to login", func(t *testing.T) {
		data := auth.LoginUserPayload{
			Email:    "hwteh1997@hotmail.com",
			Password: "123456",
		}

		// Marshal the data into JSON bytes
		jsonData, err := json.Marshal(data)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := test_utils.ExecuteRequest(req, mux)
		test_utils.CheckResponseCode(t, http.StatusOK, rr.Code)
	})

	t.Run("should not able to login", func(t *testing.T) {
		data := auth.LoginUserPayload{
			Email:    "hwteh1997@hotmail.com",
			Password: "1234567",
		}

		// Marshal the data into JSON bytes
		jsonData, err := json.Marshal(data)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := test_utils.ExecuteRequest(req, mux)
		test_utils.CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should be able to register account", func(t *testing.T) {
		data := auth.RegisterUserPayload{
			Username: "imGoot",
			Email:    "goot@gmail.com",
			Password: "veryGootPass",
		}

		// Marshal the data into JSON bytes
		jsonData, err := json.Marshal(data)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/v1/auth/user", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := test_utils.ExecuteRequest(req, mux)
		test_utils.CheckResponseCode(t, http.StatusCreated, rr.Code)
	})
	t.Run("should not be able to register account", func(t *testing.T) {
		data := auth.RegisterUserPayload{
			Username: "MehNohNah",
			Email:    "goot@gmail.com",
			Password: "veryGootPass",
		}

		// Marshal the data into JSON bytes
		jsonData, err := json.Marshal(data)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/v1/auth/user", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := test_utils.ExecuteRequest(req, mux)
		test_utils.CheckResponseCode(t, http.StatusConflict, rr.Code)
	})
	t.Run("should not be able to register account", func(t *testing.T) {
		data := auth.RegisterUserPayload{
			Username: "imGoot",
			Email:    "hwteh1997@gmail.com",
			Password: "veryGootPass",
		}

		// Marshal the data into JSON bytes
		jsonData, err := json.Marshal(data)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/v1/auth/user", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := test_utils.ExecuteRequest(req, mux)
		test_utils.CheckResponseCode(t, http.StatusConflict, rr.Code)
	})
}
