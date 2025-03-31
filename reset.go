package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetTable(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("database has been reset successfully\n")
	return nil
}
