package db

import (
	"database/sql"
	"fmt"
	"log"
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

var db *sql.DB

func Tasks(limit int) ([]*Task, error) {

	var tasks []*Task

	rows, err := db.Query("SELECT * FROM scheduler WHERE date = :date LIMIT :limit", sql.Named("limlt", limit))
	if err != nil {
		return nil, fmt.Errorf("ошибка select-запроса %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("ошибка %w", err)
		}
		tasks = append(tasks, &task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка %w", err)
	}
	return tasks, nil
}

func AddTask(task *Task) (int64, error) {
	result, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))

	if err != nil {
		log.Println(err)
		return 0, fmt.Errorf("ошибка sql-запроса")
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return 0, fmt.Errorf("ошибка при получении id")
	}
	return id, nil
}

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
