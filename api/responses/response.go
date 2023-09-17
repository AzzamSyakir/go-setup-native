package responses

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse mengembalikan respons JSON berupa pesan kesalahan.
func ErrorResponse(w http.ResponseWriter, message string, status int) {
	errorData := map[string]string{"error": message}
	response, _ := json.Marshal(errorData)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func SuccessResponse(w http.ResponseWriter, message string, data interface{}, status int) {
	type Response struct {
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

	response := Response{
		Message: message,
		Data:    data,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		ErrorResponse(w, "Gagal membuat respons JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(responseJSON)
}
