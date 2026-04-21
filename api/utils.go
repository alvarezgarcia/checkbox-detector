package main

import (
	"encoding/json"
	"net/http"
)

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func assert(cond bool, w http.ResponseWriter, msg string, code int) bool {
	if !cond {
		jsonError(w, msg, code)
	}
	return cond
}
