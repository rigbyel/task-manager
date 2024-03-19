package validate

import (
	"errors"
	"strconv"
)

var (
	ErrEmptyID   = errors.New("empty id")
	ErrInvalidID = errors.New("invalid id")
)

func ValidateID(idStr string) (int64, error) {
	if idStr == "" {
		return 0, ErrEmptyID
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, ErrInvalidID
	}

	return id, nil
}
