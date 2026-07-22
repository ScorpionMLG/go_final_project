package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func GetTask(id string) (*Task, error) {
	var task Task
	err := db.QueryRow("SELECT * FROM scheduler WHERE :id = id", sql.Named("id", id)).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return &task, err
	}
	return &task, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler
SET date = :date,
title = :title,
comment = :comment,
repeat = :repeat
WHERE id = :id;
`
	res, err := db.Exec(query, sql.Named("date", task.Date), sql.Named("title", task.Title), sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat), sql.Named("id", task.ID))
	if err != nil {
		return err
	}
	// метод RowsAffected() возвращает количество записей к которым
	// была применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func DeleteTask(id string) error {
	query := "DELETE FROM scheduler WHERE id = ?"
	res, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for deleting task`)
	}
	return nil
}

func UpdateTaskDate(task *Task) error {
	query := `UPDATE scheduler
SET date = :date
WHERE id = :id;
`
	res, err := db.Exec(query, sql.Named("date", task.Date), sql.Named("id", task.ID))
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task date`)
	}
	return nil
}

func Tasks(limit int, srch string) ([]*Task, error) {
	var query string
	tasks := make([]*Task, 0, limit)
	var rows *sql.Rows
	var err error
	if srch != "" {
		srch = strings.ToLower(srch)
		query = "SELECT * FROM scheduler WHERE rus_lower(title) LIKE :search OR rus_lower(comment) LIKE :search ORDER BY date LIMIT :limit"
		date, errParse := time.Parse("02.01.2006", srch)
		if errParse == nil {
			query = "SELECT * FROM scheduler WHERE date = :search ORDER BY date LIMIT :limit"
			dateStr := date.Format("20060102")
			rows, err = db.Query(query, sql.Named("search", dateStr), sql.Named("limit", limit))
		} else {
			srch = fmt.Sprintf("%%%s%%", srch)
			rows, err = db.Query(query, sql.Named("search", srch), sql.Named("limit", limit))
		}
	} else {
		query = "SELECT * FROM scheduler ORDER BY date LIMIT :limit"
		rows, err = db.Query(query, sql.Named("limit", limit))
	}
	if err != nil {
		return tasks, err
	}
	defer rows.Close()
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, &task)
	}
	if err = rows.Err(); err != nil {
		return tasks, err
	}
	return tasks, nil
}
