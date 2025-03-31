package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/amr-as90/rsagg/internal/database"
	"github.com/amr-as90/rsagg/internal/rss"
)

func handlerAgg(s *state, cmd command) error {
	//Check if there is a time interval argument
	if len(cmd.arguments) < 1 {
		return errors.New("not enough arguments, please provide a time between requests")
	}

	//Set the time interval to the user defined interval
	timeInterval, err := time.ParseDuration(cmd.arguments[0])
	if err != nil {
		fmt.Println("Error parsing time duration:")
		return err
	}
	//Create a time ticker and fetch the requested feed every {timeInterval}
	fmt.Printf("Collecting feeds every %s\n", timeInterval)
	ticker := time.NewTicker(timeInterval)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("Couldn't get next feeds to fetch", err)
		return
	}
	log.Println("Found a feed to fetch!")
	scrapeFeed(s.db, feed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}

	feedData, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}
	for _, item := range feedData.Channel.Item {
		fmt.Printf("Found post: %s\n", item.Title)
	}
}

func handlerBrowse(s *state, cmd command, currentUser database.User) error {
	//Check if there is a limit parameter, otherwise default to 2
	limit := 2
	if len(cmd.arguments) > 0 && cmd.arguments[0] != "" {
		parsedLimit, err := strconv.Atoi(cmd.arguments[0])
		if err != nil {
			return err
		}
		limit = parsedLimit
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: currentUser.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return err
	}

	fmt.Println("Got posts successfully!")
	fmt.Printf("Showing posts for user: %s, with ID: %s\n", currentUser.Name, currentUser.ID)
	fmt.Printf("Number of posts: %d\n", len(posts))

	for _, post := range posts {
		fmt.Println("----------------------")
		fmt.Printf("Post title: %v\n", post.Title)
		fmt.Printf("Post description: %v\n", post.Description)
		fmt.Printf("----------------------\n\n")
	}
	return nil
}
