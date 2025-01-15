package test_utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func MarshalJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

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
