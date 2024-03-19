package storage

import (
	"fmt"
)

type Manager struct {
	storage *Storage
}

func (m *Manager) AcceptQuest(userId, questId int64) error {
	op := "storage.userrepository.AcceptQuest"

	userExist, err := m.CheckUserExist(userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if !userExist {
		return fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}

	questExist, err := m.CheckQuestExist(questId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if !questExist {
		return fmt.Errorf("%s: %w", op, ErrQuestNotFound)
	}

	questCompletion, err := m.CheckQuestCompletion(userId, questId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if questCompletion {
		return fmt.Errorf("%s: %w", op, ErrQuestCompleted)
	}

	err = m.updateBalance(userId, questId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := m.storage.db.Prepare(
		"INSERT INTO quest_stat (userID, questID) VALUES ($1, $2)",
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(userId, questId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (m *Manager) updateBalance(userId, questId int64) error {
	const op = "storage.manager.updateBalance"
	balance, err := m.storage.userRepository.Balance(userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	cost, err := m.storage.questRepository.Cost(questId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	newBalance := cost + balance

	stmt, err := m.storage.db.Prepare(
		"UPDATE user SET balance = $1 WHERE id = $2",
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(newBalance, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (m *Manager) CheckUserExist(userId int64) (bool, error) {
	op := "storage.manager.CheckUserExist"

	row := m.storage.db.QueryRow(
		"SELECT COUNT(id) > 0 FROM user WHERE id = $1",
		userId,
	)

	var existFlag bool
	if err := row.Scan(&existFlag); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return existFlag, nil
}

func (m *Manager) CheckQuestExist(questId int64) (bool, error) {
	op := "storage.manager.CheckQuestExist"

	row := m.storage.db.QueryRow(
		"SELECT COUNT(id) > 0 FROM quest WHERE id = $1",
		questId,
	)

	var existFlag bool
	if err := row.Scan(&existFlag); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return existFlag, nil
}

func (m *Manager) CheckQuestCompletion(userID, questID int64) (bool, error) {
	op := "storage.manager.CheckQuestCompletion"

	row := m.storage.db.QueryRow(
		"SELECT count(id) > 0 FROM quest_stat WHERE userID = $1 AND questID = $2",
		userID,
		questID,
	)

	var existFlag bool
	if err := row.Scan(&existFlag); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return existFlag, nil
}
