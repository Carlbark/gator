package main

import (
	"context"
	"fmt"
	"strconv"
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

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := int32(2)
	if len(cmd.args) == 1 {
		num, err := strconv.Atoi(cmd.args[0])
		if err != nil || num < 1 {
			fmt.Println("Faulty limit, defaulting to 2")
		} else {
			limit = int32(num)
		}
	}
	reqData := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	}
	posts, err := s.db.GetPostsForUser(context.Background(), reqData)
	if err != nil {
		return fmt.Errorf("Failed to get posts for current user %v: %w\n", user.Name, err)
	}
	for _, post := range posts {
		fmt.Printf("Title: %s\n", post.Title)
		if !post.PublishedAt.Valid {
			fmt.Println("Published: Unknown")
		} else {
			fmt.Printf("Published: %v\n", post.PublishedAt.Time)
		}
		fmt.Printf("URL: %s\n", post.Url)
		if post.Description.Valid {
			desc := post.Description.String
			if len(desc) > 100 {
				desc = desc[:100] + "..."
			}
			fmt.Println(desc)
		}
		fmt.Println("---")
	}
	return nil
}
