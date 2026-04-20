package main

import (
	"context"
	"fmt"
	"time"

	"github.com/carlbark/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("username is required")
	}
	username := cmd.args[0]
	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("User %v does not exist in database: %w", username, err)
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}
	fmt.Printf("Logged in as: %+v\n", user)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("username is required")
	}
	username := cmd.args[0]

	regData := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	dbUser, err := s.db.CreateUser(context.Background(), regData)
	if err != nil {
		return fmt.Errorf("Failed to add user (probably already exists): %w", err)
	}
	fmt.Printf("User added: %+v\n", dbUser)
	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}
	fmt.Printf("User set to: %s\n", username)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to clear table users: %w\n", err)
	}
	fmt.Printf("Table 'users' cleared.\n")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	items, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get users from table: %w\n", err)
	}
	currentUser := s.cfg.User
	for _, item := range items {
		fmt.Printf("* %v", item.Name)
		if item.Name == currentUser {
			fmt.Printf(" (current)\n")
		} else {
			fmt.Println()
		}
	}
	return nil
}
