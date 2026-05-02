package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/carlbark/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("Time between requests (1s, 1m, 1h etc) is required")
	}
	timeBetweenReq, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("Could not parse time request: %w", err)
	}
	fmt.Printf("Collecting feeds every %v\n", timeBetweenReq)
	ticker := time.NewTicker(timeBetweenReq)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) error {
	nextToFetch, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("No feeds in database: %w\n", err)
	}
	err = s.db.MarkFeedFetched(context.Background(), nextToFetch.ID)

	if err != nil {
		return fmt.Errorf("Could not mark feed %v as fetched in database: %w\n", nextToFetch.Name, err)
	}
	rssfeed, err := fetchFeed(context.Background(), nextToFetch.Url)
	if err != nil {
		return fmt.Errorf("Failed to fetch feed %v : %w", nextToFetch.Name, err)
	}
	for _, item := range rssfeed.Channel.Item {
		if item.Title == "" || item.Link == "" {
			continue
		}
		var publishedAt sql.NullTime
		if item.PubDate != "" {
			parsedTime, err := time.Parse(time.RFC1123Z, item.PubDate)
			if err == nil {
				publishedAt = sql.NullTime{
					Time:  parsedTime,
					Valid: true,
				}
			} else {
				parsedTime, err := time.Parse(time.RFC1123, item.PubDate)
				if err == nil {
					publishedAt = sql.NullTime{
						Time:  parsedTime,
						Valid: true,
					}
				}
			}
		}
		regData := database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  item.Description != "",
			},
			PublishedAt: publishedAt,
			FeedID:      nextToFetch.ID,
		}

		err = s.db.CreatePost(context.Background(), regData)
		if err != nil {
			log.Printf("Failed inserting post: %q into database: %v", item.Title, err)
			continue
		}
	}

	return nil
}
