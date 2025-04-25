package main

import (
	// "api"
	"db"
	"server"
)

func main() {
	// api.Init()
	server.PushServer()
	db.Init("scheduler.db")
}
