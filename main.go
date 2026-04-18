package main

import (
	"fmt"
	"os"

	"github.com/carlbark/gator/internal/config"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}

	s := &state{cfg: &cfg}
	cs := &commands{cmds: make(map[string]func(*state, command) error)}
	cs.register("login", handlerLogin)

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
