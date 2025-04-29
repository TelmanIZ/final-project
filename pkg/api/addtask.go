package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/TelmanIZ/final-project/pkg/db"

	_ "modernc.org/sqlite"
)

//  функция проверяет необходимые условия при добавлении задачи, а именно:
// корректность формата даты и наличие поля title. Возвращает дату в формате int и ошибку

func checkDate(task *db.Task) (int, error) {

	if task.Title == "" {
		return 0, fmt.Errorf("не указан заголовок задачи")
	}

	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(timeFormat)
	}
	_, err := time.Parse(timeFormat, task.Date)
	if err != nil {
		return 0, fmt.Errorf("дата указана в неправильном формате")
	}

	if task.Date < now.Format(timeFormat) {
		if task.Repeat == "" {
			task.Date = now.Format(timeFormat)
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return 0, fmt.Errorf("правило повторения указано в неправильном формате")
			}
			task.Date = next
		}
	}
	dateInt, err := strconv.Atoi(task.Date)
	if err != nil {
		return 0, fmt.Errorf("не удалось конвертировать дату в число: %w", err)
	}
	return dateInt, nil
}

// NextDayHandler принимает запрос
// и возвращает следующую дату задачи или ошибку

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

// обработчик POST-запроса /api/task
// возвращает JSON объект с id

func AddTaskHandler(storage *db.DB) http.HandlerFunc {
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

		dateInt, err := checkDate(&task)
		if err != nil {
			log.Println(err)
			SendErrorResponse(w, "ошибка при проверке даты", http.StatusBadRequest)
			return
		}

		id, err := storage.AddTask(dateInt, &task)
		if err != nil {
			log.Println(err)
			SendErrorResponse(w, "ошибка при добавлении задачи", http.StatusBadRequest)
			return
		}

		response := JSONObject{ID: fmt.Sprintf("%d", id)}

		resp, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}
