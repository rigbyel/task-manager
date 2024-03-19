package storage_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/rigbyel/task-manager/internal/model"
	"github.com/rigbyel/task-manager/internal/storage"
)

func TestManager_AcceptTask(t *testing.T) {
	s, teardown := storage.TestStore(t, databaseURL)
	defer teardown("quest_stat", "quest", "user")

	path, _ := os.Executable()
	fmt.Println(path)

	q, _ := s.Quest().CreateQuest(&model.Quest{
		Name: "Go to gym",
		Cost: 88,
	})

	u, _ := s.User().CreateUser(&model.User{
		Name: "Anton Solntsev",
	})

	err := s.Manager().AcceptQuest(u.Id, q.Id)
	if err != nil {
		t.Errorf("cannot evaluate quest: %s", err.Error())
	}

	balance, _ := s.User().Balance(u.Id)

	if balance != u.Balance+q.Cost {
		t.Errorf("wrong new balance: expected %d, got %d", u.Balance+q.Cost, balance)
	}

}

func TestManager_checkQuestCompletion(t *testing.T) {
	s, teardown := storage.TestStore(t, databaseURL)
	defer teardown("quest_stat", "quest", "user")

	q, _ := s.Quest().CreateQuest(&model.Quest{
		Name: "Visit VK office",
		Cost: 88,
	})

	u, _ := s.User().CreateUser(&model.User{
		Name: "Marie Hope",
	})

	checkFlag, err := s.Manager().CheckQuestCompletion(u.Id, q.Id)
	if err != nil {
		t.Errorf("cannot check quest: %s", err.Error())
	}

	if checkFlag {
		t.Error("wrong check result: expected false, got true")
	}

	_ = s.Manager().AcceptQuest(u.Id, q.Id)

	checkFlag, err = s.Manager().CheckQuestCompletion(u.Id, q.Id)
	if err != nil {
		t.Errorf("cannot check quest: %s", err.Error())
	}

	if !checkFlag {
		t.Error("wrong check result expected true, got false")
	}

}

// TODO: implement TestManager_calcNewBalance
