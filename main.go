package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/amr-as90/rsagg/internal/config"
	"github.com/amr-as90/rsagg/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfgStruct, err := config.Read()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfgStruct.DB_URL)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	currentState := state{
		cfg: &cfgStruct,
		db:  dbQueries,
	}

	commands := commands{
		cmdNames: make(map[string]func(*state, command) error),
	}

	currentCommand := command{
		name:      os.Args[1],
		arguments: os.Args[2:],
	}

	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerUsers)
	commands.register("agg", handlerAgg)
	commands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.register("feeds", handlerGetAllFeeds)
	commands.register("follow", middlewareLoggedIn(handlerCreateFeedFollow))
	commands.register("following", middlewareLoggedIn(handlerFollows))
	commands.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	commands.register("browse", middlewareLoggedIn(handlerBrowse))

	err = commands.run(&currentState, currentCommand)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}

}
