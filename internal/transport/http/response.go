package http

import (
	"encoding/json"
	nethttp "net/http"
)

type envelope struct {
	Success bool    `json:"success"`
	Data    any     `json:"data"`
	Error   *string `json:"error"`
}

func writeOK(w nethttp.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(envelope{Success: true, Data: data, Error: nil})
}

func writeError(w nethttp.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(envelope{Success: false, Data: nil, Error: &msg})
}
