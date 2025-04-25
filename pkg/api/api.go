package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	// "github.com/go-chi/chi"
)

type JSONObject struct {
	ID    string `json:"id,omitempty"`
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

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

func NextDayHandler(res http.ResponseWriter, rep *http.Request) {
	nowStr := rep.FormValue("now")
	dateStr := rep.FormValue("date")
	repeatStr := rep.FormValue("repeat")

	now, err := time.Parse(timeFormat, nowStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	taskDay, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	res.Write([]byte(taskDay))
}

// func writeJSON(w http.ResponseWriter, data any) {
// 	resp, err := json.Marshal(data)
// 	if err != nil {
// 		log.Println(err)
// 		http.Error(w, "Ошибка сериализации", http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
// 	w.WriteHeader(http.StatusOK)
// 	_, err = w.Write(resp)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

// // func taskHandler(w http.ResponseWriter, r *http.Request) {
// // 	switch r.Method {
// // 	case http.MethodPost:
// // 		addTaskHandler(w, r)
// 		// case http.MethodGet:
// 		// 	getTaskHandler(w, r)
// 		// case http.MethodPut:
// 		// 	updateTaskHandler(w, r)
// 	}

// }

// func Init() {
// 	r := chi.NewRouter()
// 	r.Get("/api/nextdate", nextDayHandler)
// 	r.Post("/api/task", addTaskHandler)
// 	// 	http.HandleFunc("/api/tasks", tasksHandler)
// 	// 	http.HandleFunc("/api/task/done", doneTaskHandler)
// }
