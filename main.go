package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/carlbark/gator/internal/config"
	"github.com/carlbark/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}

	db, err := sql.Open("postgres", cfg.Url)
	if err != nil {
		fmt.Println("Error opening database connection:", err)
		return
	}
	dbQueries := database.New(db)

	s := &state{cfg: &cfg, db: dbQueries}
	cs := &commands{cmds: make(map[string]func(*state, command) error)}
	cs.register("login", handlerLogin)
	cs.register("register", handlerRegister)
	cs.register("reset", handlerReset)
	cs.register("users", handlerUsers)

	input := os.Args
	if len(input) < 2 {
		fmt.Println("Command is required")
		os.Exit(1)
	}
	cmd := command{name: input[1], args: input[2:]}
	err = cs.run(s, cmd)
	if err != nil {
		fmt.Println("Error running command:", err)
		os.Exit(1)
	}

}
