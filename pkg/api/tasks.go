package api

import (
	"bytes"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/TelmanIZ/final-project/pkg/db"
	// "fmt"
	"net/http"
	// "time"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func TasksHandler(storage *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks, err := storage.Tasks(50)
		if err != nil {
			SendErrorResponse(w, "ошибка при отображении задачи", http.StatusBadRequest)
			return
		}
		if tasks == nil {
			response := TasksResp{Tasks: []*db.Task{}}
			writeJSON(w, response)
			return
		}
		writeJSON(w, TasksResp{
			Tasks: tasks,
		})
	}
}

func GetTaskHandler(storage *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		resp, err := storage.GetTask(id)
		if err != nil {
			SendErrorResponse(w, "не указан идентификатор", http.StatusInternalServerError)
			return
		}
		writeJSON(w, resp)
	}
}

func UpdateTaskHandler(storage *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var task db.Task
		var buf bytes.Buffer

		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(buf.Bytes(), &task)
		if err != nil {
			log.Println(err)
			SendErrorResponse(w, "ошибка десериализации JSON", http.StatusBadRequest)
			return
		}

		err = storage.GetForUpdateTask(task)
		if err != nil {
			log.Println(err)
			SendErrorResponse(w, "ошибка получения задачи", http.StatusBadRequest)
			return
		}

		dateInt, err := checkDate(&task)
		if err != nil {
			log.Println(err)
			SendErrorResponse(w, "ошибка при проверке даты", http.StatusBadRequest)
			return
		}

		err = storage.UpdateTask(dateInt, task)
		if err != nil {
			http.Error(w, "ошибка при редактировании задачи", http.StatusInternalServerError)
		}
		writeJSON(w, task)
	}
}

func DoneTaskHandler(storage *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var s struct{}

		id := r.URL.Query().Get("id")
		if id == "" {
			SendErrorResponse(w, "не указан идентификатор", http.StatusBadRequest)
			return
		}

		task, err := storage.GetTask(id)
		if err != nil {
			SendErrorResponse(w, "ошибка при выгрузке задачи", http.StatusInternalServerError)
			return
		}

		if task.Repeat == "" {
			err = storage.DeleteTask(task.ID)
			if err != nil {
				SendErrorResponse(w, "ошибка при удалении", http.StatusInternalServerError)
				return
			}
			writeJSON(w, s)
			return
		}

		now := time.Now()
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			SendErrorResponse(w, "ошибка обработки даты", http.StatusInternalServerError)
			return
		}
		dateInt, err := strconv.Atoi(next)
		if err != nil {
			log.Println(err)
			SendErrorResponse(w, "ошибка при преобразовании строки в число", http.StatusBadRequest)
			return
		}
		err = storage.UpdateDate(dateInt, task.ID)
		if err != nil {
			SendErrorResponse(w, "ошибка при добавлении даты", http.StatusInternalServerError)
			return
		}
		task.Date = next
		writeJSON(w, s)
	}
}

func DeleteTaskHandler(storage *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var s struct{}

		id := r.URL.Query().Get("id")

		err := storage.DeleteTask(id)
		if err != nil {
			SendErrorResponse(w, "ошибка при удалении", http.StatusBadRequest)
			return
		}
		writeJSON(w, s)
	}
}
