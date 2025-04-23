package api

import (
	"database/sql"
	"db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func checkDate(task *db.Task) error {
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(timeFormat)
	}
	_, err := time.Parse(timeFormat, task.Date)
	if err != nil {
		return fmt.Errorf("дата указана в неправильном формате")
	}

	if task.Date < now.Format(timeFormat) {
		if task.Repeat == "" {
			task.Date = now.Format(timeFormat)
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("правило повторения указано в неправильном формате")
			}
			task.Date = next
		}
	}
	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {

	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		errTxt := fmt.Errorf("ошибка десериализации JSON")
		errorResp := JSONObject{Error: errTxt}
		writeJSON(w, errorResp)
		return
	}

	if task.Title == "" {
		errTxt := fmt.Errorf("не указан заголовок задачи")
		errorResp := JSONObject{Error: errTxt}
		writeJSON(w, errorResp)
		return
	}

	err = checkDate(&task)
	if err != nil {
		errorResp := JSONObject{Error: err}
		writeJSON(w, errorResp)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%w", err), http.StatusBadRequest)
	}

	res := JSONObject{ID: fmt.Sprintf("%d", id)}
	writeJSON(w, res)
}
