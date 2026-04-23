package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to clear table users: %w\n", err)
	}
	fmt.Printf("Table 'users' cleared.\n")
	return nil
}
