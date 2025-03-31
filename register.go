package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/amr-as90/rsagg/internal/database"
	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("no arguments found, a username is required")
	}

	user, _ := s.db.GetUser(context.Background(), cmd.arguments[0])
	if user.Name == cmd.arguments[0] {
		fmt.Println("username already exists.")
	}

	userID, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	user, err = s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
	})
	if err != nil {
		fmt.Println("unable to create new user in the database")
		os.Exit(1)
	}

	err = s.cfg.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}

	fmt.Printf("New user %s has been set!\n", user.Name)
	return nil
}
