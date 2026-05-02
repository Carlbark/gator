package main

import (
	"context"
	"fmt"

	"github.com/carlbark/gator/internal/config"
	"github.com/carlbark/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	if f, exists := c.cmds[cmd.name]; exists {
		return f(s, cmd)
	}
	return fmt.Errorf("command not found: %s", cmd.name)
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		currentUser := s.cfg.User
		user, err := s.db.GetUser(context.Background(), currentUser)
		if err != nil {
			return fmt.Errorf("User %v does not exist in database: %w", currentUser, err)
		}
		return handler(s, cmd, user)
	}
}
