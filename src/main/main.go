package main

import (
	"encoding/json"
	"fmt"

	"Flow/src/lexer"
	"Flow/src/parser"
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

	p := parser.New(lexer.New(`fn (a, b) { c = a * b + 2; q = false; return c; }`))
	program := p.ParseProgram()
	fmt.Println(program.String())
	data, _ := json.Marshal(program)
	fmt.Println(string(data))
}
