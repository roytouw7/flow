package main

import (
	"fmt"
)

var (
	letters = []rune{'a', 'b', 'c', 'd', 'e'}
)

func main() {
	//user, err := user2.Current()
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Printf("Hello %s! this is the Flow programming language!\n", user.Username)
	//fmt.Printf("Feel free to enter any commands!\n")
	//fmt.Printf("Enter .flow file to read it as input instead!\n")
	//repl.Start(os.Stdin, os.Stdout)

	fmt.Println(string(letters[0:1]))
}

func hasNext(n int) bool {
	return n-1 < len(letters)
}
