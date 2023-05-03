package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/GzArthur/interpreter/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is a simple interpreter!\n",
		user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.StartREPL(os.Stdin, os.Stdout)
}
