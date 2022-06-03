package utils

import (
	"encoding/json"
	"net/http"
)

func CorsHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
}

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	CorsHeader(w)
	json.NewEncoder(w).Encode(data)
}
