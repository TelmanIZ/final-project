package api

// import (
// 	"db"
// 	"fmt"
// 	"net/http"
// 	"time"
// )

// type TasksResp struct {
// 	Tasks []*db.Task `json:"tasks"`
// }

// func tasksHandler(w http.ResponseWriter, r *http.Request) {
// 	tasks, err := db.Tasks(50)
// 	if err != nil {
// 		errorResp := JSONObject{Error: fmt.Errorf("ошибка при отображении задачи: %w", err)}
// 		writeJSON(w, errorResp)
// 		return
// 	}
// 	res := TasksResp{Tasks: tasks}
// 	writeJSON(w, res)
// 	if tasks == nil {
// 		response := TasksResp{Tasks: []*db.Task{}}
// 		writeJSON(w, response)
// 		return
// 	}

// }

// func getTaskHandler(w http.ResponseWriter, r *http.Request) {
// 	id := r.URL.Query().Get("id")
// 	resp, err := db.GetTask(id)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	writeJSON(w, resp)
// }

// func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
// 	var task *db.Task
// 	err := db.UpdateTask(task)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// 	writeJSON(w, task)
// }

// func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
// 	id := r.URL.Query().Get("id")
// 	task, err := db.GetTask(id)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	if task.Repeat == "" {
// 	err = db.DeleteTask(task.ID)
// 	if err != nil {
// 		errorResp := JSONObject{Error: fmt.Errorf("ошибка при удалении задачи: %w", err)}
// 		writeJSON(w, errorResp)
// 		return
// 	}
// 	writeJSON(w, struct{}{})
// }
// 		now := time.Now()
// 		next, err := NextDate(now, task.Date, task.Repeat)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		err = db.UpdateDate(next, task.ID)
// 		if err != nil {
// 			errorResp := JSONObject{Error: fmt.Errorf("ошибка при добавлении задачи: %w", err)}
// 			writeJSON(w, errorResp)
// 			return
// // 	}
// // }
