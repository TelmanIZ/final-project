package api

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSONObject используется при формиравании JSON-объекта
type JSONObject struct {
	ID    string `json:"id,omitempty"`
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

// SendErrorResponse отправляет текст ошибки в формате JSON

func SendErrorResponse(res http.ResponseWriter, errorMessage string, statusCode int) {
	response := JSONObject{Error: errorMessage}

	resp, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		http.Error(res, "ошибка при сериализации JSON", http.StatusInternalServerError)
	}

	res.Header().Set("Content-Type", "application/json;charset=UTF-8")
	res.WriteHeader(statusCode)
	_, err = res.Write(resp)
	if err != nil {
		log.Println(err)
	}
}

// writeJSON отправляет interface в формате JSON

func writeJSON(w http.ResponseWriter, data any) {
	resp, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Ошибка сериализации", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		log.Println(err)
	}
}
