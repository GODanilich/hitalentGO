package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"GODanilich/hitalentGO/internal/apperr"
)

const maxBodyBytes = 1 << 20 // 1MB

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {

		return apperr.Validation("invalid JSON body")
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return apperr.Validation("invalid JSON body")
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeAppError(w http.ResponseWriter, err error) {
	if ae, ok := apperr.As(err); ok {
		writeJSON(w, ae.HTTPStatus, map[string]any{
			"error": map[string]any{
				"code":    ae.Code,
				"message": ae.Message,
			},
		})
		return
	}

	writeJSON(w, http.StatusInternalServerError, map[string]any{
		"error": map[string]any{
			"code":    apperr.CodeInternal,
			"message": "internal error",
		},
	})
}
