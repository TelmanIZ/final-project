package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

const schema string = `CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date INTEGER NOT NULL DEFAULT 20060102,
	title TEXT NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat TEXT CHECK (LENGTH(repeat) <= 128) NOT NULL DEFAULT "");
	CREATE INDEX scheduler_date ON scheduler (date);`

func Init(dbFile string) error {

	_, err := os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	if install {
		_, err := db.Exec(schema)
		if err != nil {
			return err
		}
	}
	return nil
}
