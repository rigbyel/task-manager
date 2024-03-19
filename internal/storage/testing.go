package storage

import (
	"fmt"
	"testing"
)

// creating storage for testing
func TestStore(t *testing.T, databaseURL string) (*Storage, func(...string)) {
	t.Helper()

	s, err := New(databaseURL)
	if err != nil {
		t.Fatal(err)
	}

	return s, func(tables ...string) {
		if len(tables) > 0 {
			for _, tbl := range tables {
				if _, err := s.db.Exec(fmt.Sprintf("DELETE FROM %s", tbl)); err != nil {
					t.Fatal(err)
				}
			}
		}

		s.Stop()
	}
}
