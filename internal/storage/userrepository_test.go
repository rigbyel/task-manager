package storage_test

import (
	"testing"

	"github.com/rigbyel/task-manager/internal/model"
	"github.com/rigbyel/task-manager/internal/storage"
)

func TestUserRepository_CreateUser(t *testing.T) {
	s, teardown := storage.TestStore(t, databaseURL)
	defer teardown("user")

	u, err := s.User().CreateUser(&model.User{
		Name: "Mary Lucky",
	})
	if err != nil {
		t.Errorf("cannot create user: %s", err.Error())
	}

	if u == nil {
		t.Error("empty user")
	}
}

func TestUserRepository_Balance(t *testing.T) {
	s, teardown := storage.TestStore(t, databaseURL)
	defer teardown("user")

	inputBalance := 120
	u, _ := s.User().CreateUser(&model.User{
		Name:    "Jane Happy",
		Balance: inputBalance,
	})

	outputBalance, err := s.User().Balance(u.Id)
	if err != nil {
		t.Errorf("cannot get user balance: %s", err.Error())
	}

	if inputBalance != outputBalance {
		t.Errorf("wrong balance: expected %d, got %d", inputBalance, outputBalance)
	}
}

func TestUserRepository_History(t *testing.T) {
	s, teardown := storage.TestStore(t, databaseURL)
	defer teardown("quest_stat", "user", "quest")

	u, _ := s.User().CreateUser(&model.User{Name: "James Bond"})

	quests := []*model.Quest{
		{Name: "quest1", Cost: 1},
		{Name: "quest2", Cost: 1},
		{Name: "quest3", Cost: 1},
	}

	for _, q := range quests {
		q, _ = s.Quest().CreateQuest(q)
		_ = s.Manager().AcceptQuest(u.Id, q.Id)
	}

	history, err := s.User().History(u.Id)
	if err != nil {
		t.Error("cannot get user history: ", err)
	}

	if len(history) != len(quests) {
		t.Errorf("wrong history length: expected %d, got %d", len(quests), len(history))
	}
}
