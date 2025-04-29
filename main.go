package main

import (
	"github.com/TelmanIZ/final-project/pkg/db"
	"github.com/TelmanIZ/final-project/pkg/server"
)

func main() {
	db, err := db.Init("./pkg/db/scheduler.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	server.PushServer(db)
}
