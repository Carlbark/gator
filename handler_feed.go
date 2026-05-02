package main

import (
	"context"
	"fmt"
	"time"

	"github.com/carlbark/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("Feed name and URL are required")
	}
	feedname := cmd.args[0]
	url := cmd.args[1]

	regData := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedname,
		Url:       url,
		UserID:    user.ID,
	}

	dbFeed, err := s.db.CreateFeed(context.Background(), regData)
	if err != nil {
		return fmt.Errorf("Failed to add feed (probably already exists): %w", err)
	}

	regFeedFollowData := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    dbFeed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), regFeedFollowData)
	if err != nil {
		return fmt.Errorf("Failed to follow feed (you are probably already following this feed): %w", err)
	}

	fmt.Printf("Feed added: %+v\n", dbFeed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	items, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get feeds from table: %w\n", err)
	}
	for _, item := range items {
		fmt.Printf("%-6s %v\n", "Feed:", item.Name)
		fmt.Printf("%-6s %v\n", "URL:", item.Url)
		fmt.Printf("%-6s %v\n", "User:", item.UserName)
		fmt.Println("---")
	}

	return nil
}
