package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
	"github.com/rigbyel/task-manager/internal/model"
)

type questRepository struct {
	storage *Storage
}

func (r *questRepository) CreateQuest(q *model.Quest) (*model.Quest, error) {
	op := "storage.questrepository.CreateQuest"

	stmt, err := r.storage.db.Prepare(
		"INSERT INTO quest (name, cost) VALUES ($1, $2)",
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(q.Name, q.Cost)
	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return nil, fmt.Errorf("%s: %w", op, ErrQuestExists)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	q.Id = id

	return q, nil
}

func (r *questRepository) Cost(questID int64) (int, error) {
	op := "storage.questRepository.Cost"

	row := r.storage.db.QueryRow(
		"SELECT cost FROM quest WHERE id = $1",
		questID,
	)

	var cost int
	if err := row.Scan(&cost); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("%s: %w", op, ErrQuestNotFound)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return cost, nil
}
