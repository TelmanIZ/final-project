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
		return fmt.Errorf("ошибка при парсинге")
	}

	if task.Date < now.Format(timeFormat) {
		if task.Repeat == "" {
			task.Date = now.Format(timeFormat)
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("ошибка при парсинге даты: %w", err)
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
		log.Println(err)
		http.Error(w, "Ошибка десериализации", http.StatusBadRequest)
		return
	}
	err = checkDate(&task)
	if err != nil {
		log.Println(err)
		http.Error(w, "Неверный формат", http.StatusBadRequest)
		return
	}
	result, err := DB.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))

	if err != nil {
		log.Println(err)
		http.Error(w, "ошибка при парсинге времени", http.StatusBadRequest)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		http.Error(w, "ошибка при получении id", http.StatusBadRequest)
		return
	}
	res := db.Task{ID: fmt.Sprintf("%d", id), Date: task.Date, Title: task.Title, Comment: task.Comment, Repeat: task.Repeat}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(res)
}
