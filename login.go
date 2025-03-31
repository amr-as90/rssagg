package main

import (
	"context"
	"errors"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("no arguments found, a username is required")
	}

	user, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err != nil {
		return errors.New("user doesn't exist in the database")
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("User %s logged in successfully!\n", user.Name)
	return nil
}
