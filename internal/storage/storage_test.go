package storage_test

import (
	"database/sql"
	"os"
	"testing"
)

var (
	databaseURL = "teststorage.db"
)

func TestMain(m *testing.M) {
	db, err := sql.Open("sqlite3", databaseURL)
	if err != nil {
		panic(err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS quest (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			cost INT NOT NULL
		);`,
	)
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}

	stmt, err = db.Prepare(`
		CREATE TABLE IF NOT EXISTS user(
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			balance INTEGER
		);`,
	)
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}

	stmt, err = db.Prepare(`
		CREATE TABLE IF NOT EXISTS quest_stat(
			id INTEGER PRIMARY KEY,
			userID INTEGER NOT NULL,
			questID INTEGER NOT NULL,
			FOREIGN KEY (userID) REFERENCES user(id) ON DELETE CASCADE
		);`,
	)
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}
