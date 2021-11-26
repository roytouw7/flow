package main

import (
	"Flow/src/repl"
	"fmt"
	"os"
	user2 "os/user"
)

func main() {
	user, err := user2.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! this is the Flow programming language!\n", user.Username)
	fmt.Printf("Feel free to enter any commands!\n")
	fmt.Printf("Enter .flow file to read it as input instead!\n")
	repl.Start(os.Stdin, os.Stdout)
}
