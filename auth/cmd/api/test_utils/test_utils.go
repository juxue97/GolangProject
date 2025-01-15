package test_utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type (
	userKey       string
	targetUserKey string
)

const (
	userCtx       userKey       = "user"
	targetUserCtx targetUserKey = "targetUser"
)

func ExecuteRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

func CheckResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d", expected, actual)
	}
}

func TestRequestWithPayload(t *testing.T, method, url string, payload interface{}, expectedStatus int) *http.Request {
	t.Helper()
	// Marshal the payload into JSON bytes
	jsonData, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	return req
}

func TestRequestWithoutPayload(t *testing.T, method, url string, expectedStatus int) *http.Request {
	t.Helper()

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	return req
}
