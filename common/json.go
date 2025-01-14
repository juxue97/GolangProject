package common

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
	return nil
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 // 1mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
	// return json.NewDecoder(r.Body).Decode(data)
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}

func JsonResponse(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}

	return WriteJSON(w, status, &envelope{Data: data})
}
