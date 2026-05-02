package main

import (
	"context"
	"fmt"
	"time"

	"github.com/carlbark/gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("URL is required")
	}
	url := cmd.args[0]

	feed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Feed with URL: %v does not exist in database: %w", url, err)
	}

	regData := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	dbFeedFollow, err := s.db.CreateFeedFollow(context.Background(), regData)
	if err != nil {
		return fmt.Errorf("Failed to follow feed (you are probably already following this feed): %w", err)
	}
	fmt.Printf("Feed followed: %+v\n", dbFeedFollow.FeedName)
	fmt.Printf("By user: %+v\n", dbFeedFollow.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	items, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Failed to get feeds for the current user: %v from table: %w\n", user.Name, err)
	}
	fmt.Println("You are following these feeds:")
	for _, item := range items {
		fmt.Printf("%v\n", item.FeedName)
	}

	return nil
}

func handlerUnFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("URL is required")
	}
	url := cmd.args[0]

	feed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Feed with URL: %v does not exist in database: %w", url, err)
	}

	regData := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	err = s.db.DeleteFeedFollow(context.Background(), regData)
	if err != nil {
		return fmt.Errorf("Failed to unfollow feed: %v for the current user: %v from table: %w\n", feed.Name, user.Name, err)
	}
	fmt.Printf("You have unfollowed this feed: %v\n", feed.Name)
	return nil
}
