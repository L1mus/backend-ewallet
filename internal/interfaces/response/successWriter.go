package response

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Status string `json:"status"`
	Data   any    `json:"data,omitempty"`
	Meta   any    `json:"meta,omitempty"`
}

func WriteSuccess(w http.ResponseWriter, statusCode int, data any) {
	writeSuccessJSON(w, statusCode, SuccessResponse{
		Status: "success",
		Data:   data,
	})
}

func WriteSuccessWithMeta(w http.ResponseWriter, statusCode int, data any, meta any) {
	writeSuccessJSON(w, statusCode, SuccessResponse{
		Status: "success",
		Data:   data,
		Meta:   meta,
	})
}

func WriteNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func writeSuccessJSON(w http.ResponseWriter, statusCode int, payload SuccessResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}
