package server

import (
	"api"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func PushServer() {
	r := chi.NewRouter()
	webDir := "./web"
	r.Handle("/", http.FileServer(http.Dir(webDir)))
	r.Get("/api/nextdate", api.NextDayHandler)
	r.Post("/api/task", api.AddTaskHandler)

	if err := http.ListenAndServe(":7540", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
