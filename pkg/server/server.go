package server

import (
	"fmt"
	"net/http"

	"github.com/TelmanIZ/final-project/pkg/api"
	"github.com/TelmanIZ/final-project/pkg/db"

	"github.com/go-chi/chi"
)

func PushServer(storage *db.DB) {
	r := chi.NewRouter()
	webDir := "./web"
	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir(webDir))))

	r.Get("/api/nextdate", api.NextDayHandler)
	r.Post("/api/task", api.AddTaskHandler(storage))
	r.Get("/api/tasks", api.TasksHandler(storage))
	r.Get("/api/task", api.GetTaskHandler(storage))
	r.Put("/api/task", api.UpdateTaskHandler(storage))
	r.Post("/api/task/done", api.DoneTaskHandler(storage))
	r.Delete("/api/task", api.DeleteTaskHandler(storage))

	if err := http.ListenAndServe(":7540", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
