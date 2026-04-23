package main

import (
	"context"
	"fmt"
)

func handlerAgg(s *state, cmd command) error {
	rssfeed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("Failed to aggregate feed: %w", err)
	}
	fmt.Printf("%+v\n", rssfeed)
	return nil
}
