package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/amr-as90/rsagg/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, currentUser database.User) error {
	//Check if there are enough arguments provided with the command
	if len(cmd.arguments) < 2 {
		return errors.New("no arguments found, please provide a URL and feed name")
	}

	//Create a feed with the variables provided
	createdFeed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
		Url:       cmd.arguments[1],
		UserID:    currentUser.ID,
	})
	if err != nil {
		return fmt.Errorf("ERROR: %s", err)
	}

	//Create a feed follow for the newly created feed
	s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    createdFeed.UserID,
		FeedID:    createdFeed.ID,
	})

	fmt.Println(createdFeed)

	return nil
}

func handlerGetAllFeeds(s *state, cmd command) error {
	//Get all the feeds in the database
	feeds, err := s.db.GetAllFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("ERROR: Couldn't get feeds from database. %s", err)
	}

	for _, feed := range feeds {
		fmt.Printf("- '%s'\n", feed.Name)
		fmt.Printf("- '%s'\n", feed.Name_2)
	}

	return nil
}

func handlerCreateFeedFollow(s *state, cmd command, currentUser database.User) error {
	//Check if URL argument is available
	if len(cmd.arguments) < 1 {
		return errors.New("not enough arguments found. Please provide a URL")
	}

	//Get the feed information based on the URL
	feedInfo, err := s.db.GetFeedByURL(context.Background(), cmd.arguments[0])
	if err != nil {
		return errors.New("no feed exists with this URL")
	}

	//Make the database query to create the feed follow
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUser.ID,
		FeedID:    feedInfo.ID,
	})
	if err != nil {
		return fmt.Errorf("unable to create the feed follow: %s", err)
	}

	fmt.Printf("User %s is now following feed %s", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func handlerFollows(s *state, cmd command, currentUser database.User) error {

	feeds, err := s.db.GetAllFeedsForUser(context.Background(), currentUser.ID)
	if err != nil {
		return errors.New("unable to recover feeds for this user")
	}

	//Print feeds currently being followed by current user
	fmt.Printf("Feeds currently being followed by %s:\n", currentUser.Name)

	for _, feed := range feeds {
		fmt.Printf("%s\n", feed)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, currentUser database.User) error {
	//Check if URL argument is available
	if len(cmd.arguments) < 1 {
		return errors.New("not enough arguments found. Please provide a URL")
	}

	err := s.db.UnfollowFeed(context.Background(), database.UnfollowFeedParams{
		UserID: currentUser.ID,
		Url:    cmd.arguments[0],
	})
	if err != nil {
		return errors.New("unable to unfollow feed, are you are following this feed?")
	}
	return nil
}
