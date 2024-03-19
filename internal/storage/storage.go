package storage

import (
	"errors"
)

var (
	ErrUserExists     = errors.New("user already exists")
	ErrQuestExists    = errors.New("quest already exists")
	ErrUserNotFound   = errors.New("user not found")
	ErrQuestNotFound  = errors.New("quest not found")
	ErrQuestCompleted = errors.New("quest already completed")
	ErrAppNotFound    = errors.New("app not found")
)
