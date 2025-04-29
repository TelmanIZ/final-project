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

type DB struct {
	db *sql.DB
}

func Init(dbFile string) (*DB, error) {

	_, err := os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	storage := &DB{db: db}

	if install {
		_, err := db.Exec(schema)
		if err != nil {
			return nil, err
		}
	}
	return storage, nil
}

func (storage *DB) Close() error {
	if err := storage.db.Close(); err != nil {
		return fmt.Errorf("ошибка при закрытии базы данных: %w", err)
	}
	return nil
}

func (Storage *DB) AddTask(dateInt int, task *Task) (int64, error) {
	result, err := Storage.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", dateInt),
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

func (storage *DB) Tasks(limit int) ([]*Task, error) {

	var tasks []*Task

	rows, err := storage.db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?", limit)
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

func (storage *DB) GetTask(id string) (Task, error) {

	var task Task

	row := storage.db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return Task{}, fmt.Errorf("задача не найдена %w", err)
		}
		log.Println(err)
		return Task{}, fmt.Errorf("ошибка при выгрузке задачи %w", err)
	}
	return task, nil
}

func (storage *DB) GetTaskDone(id string) (Task, error) {
	var task Task

	rows, err := storage.db.Query("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))
	if err != nil {
		return Task{}, fmt.Errorf("ошибка при select-запросе %w", err)

	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		return Task{}, fmt.Errorf("ошибка при select-запросе %w", err)
	}
	if err := rows.Err(); err != nil {
		return Task{}, fmt.Errorf("ошибка при select-запросе %w", err)
	}
	return task, nil
}

func (storage *DB) GetForUpdateTask(task Task) error {

	row := storage.db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", task.ID))
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("ошибка при добавлении задачи %w", err)
	}
	return nil
}

func (storage *DB) UpdateTask(dateInt int, task Task) error {

	res, err := storage.db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", dateInt),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))

	if err != nil {
		return fmt.Errorf("ошибка %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func (storage *DB) DeleteTask(id string) error {
	result, err := storage.db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("ошибка %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения удаленных строк %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача с таким id не найдена %w", err)
	}
	return nil
}

func (storage *DB) UpdateDate(next int, id string) error {

	res, err := storage.db.Exec("UPDATE scheduler SET date = :date WHERE id = :id",
		sql.Named("date", next),
		sql.Named("id", id))

	if err != nil {
		return fmt.Errorf("ошибка %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}

	return nil
}
