package db

// import (
// 	"database/sql"
// 	"fmt"
// )

// func GetTask(id string) (*Task, error) {
// 	var task *Task
// 	row := db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))
// 	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
// 	if err != nil {
// 		return nil, fmt.Errorf("ошибка при sql-запросе: %w", err)
// 	}
// 	return task, nil
// }

// func UpdateTask(task *Task) error {
// 	query := `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id= :id`
// 	res, err := db.Exec(query, sql.Named("date", task.Date),
// 		sql.Named("title", task.Title),
// 		sql.Named("comment", task.Comment),
// 		sql.Named("repeat", task.Repeat),
// 		sql.Named("id", task.ID))
// 	if err != nil {
// 		return err
// 	}
// 	count, err := res.RowsAffected()
// 	if err != nil {
// 		return err
// 	}
// 	if count == 0 {
// 		return fmt.Errorf(`incorrect id for updating task`)
// 	}
// 	return nil
// }

// func DeleteTask(id string) error {
// 	_, err := db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func UpdateDate(next string, id string) error {
// 	res, err := db.Exec("UPDATE scheduleer SET date = :date WHERE id = :id",
// 		sql.Named("date", next),
// 		sql.Named(":id", id))
// 	if err != nil {
// 		return err
// 	}
// 	count, err := res.RowsAffected()
// 	if err != nil {
// 		return err
// 	}
// 	if count == 0 {
// 		return fmt.Errorf(`incorrect id for updating task`)
// 	}
// 	return nil
// }
