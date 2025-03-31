package main

import (
	"fmt"

	"github.com/amr-as90/rsagg/internal/config"
	"github.com/amr-as90/rsagg/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	cmdNames map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmdNames[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	function, exists := c.cmdNames[cmd.name]
	if !exists {
		return fmt.Errorf("command does not exist")
	}
	err := function(s, cmd)
	if err != nil {
		return err
	}
	return nil
}
