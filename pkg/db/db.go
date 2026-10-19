package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"os"
	"strings"

	"modernc.org/sqlite"
)

const schema = `
CREATE TABLE scheduler (
id INTEGER PRIMARY KEY AUTOINCREMENT,
date CHAR(8) NOT NULL DEFAULT "",
title VARCHAR(128) NOT NULL DEFAULT "",
comment TEXT DEFAULT "",
repeat VARCHAR(128) DEFAULT ""
);

CREATE INDEX idx_scheduler_date ON scheduler(date);
`

var db *sql.DB

func Init(dbFile string) error {
	env := os.Getenv("TODO_DBFILE")
	if env != "" {
		dbFile = env
	}

	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	// Регистрация функции 'rus_lower', которая корректно переводит кириллицу в нижний регистр
	err = sqlite.RegisterDeterministicScalarFunction(
		"rus_lower",
		1,
		func(ctx *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
			if args[0] == nil {
				return nil, nil
			}
			val, ok := args[0].(string)
			if !ok {
				return nil, fmt.Errorf("expected string argument")
			}
			return strings.ToLower(val), nil
		},
	)
	if err != nil {
		return fmt.Errorf("failed to register lower function: %w", err)
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("No connection with db: %w", err)
	}

	if install {
		_, err := db.Exec(schema)
		if err != nil {
			return err
		}
	}
	return nil
}

func Close() {
	db.Close()
}
