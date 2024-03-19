package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
	"github.com/rigbyel/task-manager/internal/model"
)

type userRepository struct {
	storage *Storage
}

func (r *userRepository) CreateUser(u *model.User) (*model.User, error) {
	op := "storage.userRepository.CreateUser"

	stmt, err := r.storage.db.Prepare(
		"INSERT INTO user (name, balance) VALUES ($1, $2)",
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(u.Name, u.Balance)
	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return nil, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	u.Id = id

	return u, nil
}

func (r *userRepository) History(userID int64) ([]string, error) {
	op := "storage.userRepository.History"

	rows, err := r.storage.db.Query(
		"SELECT quest.name FROM quest, quest_stat WHERE quest_stat.userID = $1 AND quest_stat.questID = quest.id;",
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var questList []string

	for rows.Next() {
		var questName string
		rows.Scan(&questName)

		questList = append(questList, questName)
	}
	rows.Close()

	return questList, nil
}

func (r *userRepository) Balance(userID int64) (int, error) {
	op := "storage.userRepository.Balance"

	row := r.storage.db.QueryRow(
		"SELECT balance FROM user WHERE id = $1",
		userID,
	)

	var balance int
	if err := row.Scan(&balance); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return balance, nil
}
