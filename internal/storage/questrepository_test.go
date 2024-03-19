package storage_test

import (
	"testing"

	"github.com/rigbyel/task-manager/internal/model"
	"github.com/rigbyel/task-manager/internal/storage"
)

func TestQuestRepository_CreateQuest(t *testing.T) {
	s, teardown := storage.TestStore(t, databaseURL)
	defer teardown("quest")

	q, err := s.Quest().CreateQuest(&model.Quest{
		Name: "Eat Well",
		Cost: 118,
	})
	if err != nil {
		t.Errorf("cannot create quest: %s", err.Error())
	}

	if q == nil {
		t.Error("empty quest")
	}
}

func TestQuestRepository_Cost(t *testing.T) {
	s, teardown := storage.TestStore(t, databaseURL)
	defer teardown("quest")

	inputCost := 98
	q, _ := s.Quest().CreateQuest(&model.Quest{
		Name: "Bake a cake",
		Cost: inputCost,
	})

	outputCost, err := s.Quest().Cost(q.Id)
	if err != nil {
		t.Errorf("cannot get quest cost: %s", err.Error())
	}

	if inputCost != outputCost {
		t.Errorf("wrong cost: expected %d, got %d", inputCost, outputCost)
	}
}
