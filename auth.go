package main

import (
	"context"
	"errors"

	"github.com/amr-as90/rsagg/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {

		//Set the user to the currently defined user
		user := s.cfg.CurrentUserName
		if user == "" {
			return errors.New("user is currently not logged in")
		}

		//Get the current user from the database and store the information in a struct
		currentUser, err := s.db.GetUser(context.Background(), user)
		if err != nil {
			return errors.New("user doesn't exist in the database")
		}
		return handler(s, cmd, currentUser)
	}
}
