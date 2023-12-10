package main

import (
	"os"

	"Flow/src/eval"
	"Flow/src/lexer"
	"Flow/src/object"
	"Flow/src/parser"
)

func main() {
	//user, err := user.Current()
	//if err != nil {
	//	panic(err)
	//}

	//fmt.Printf("Hello %s! this is the Flow programming language!\n", user.Username)
	//fmt.Printf("Feel free to enter any commands!\n")
	//fmt.Printf("Enter .flow file to read it as input instead!\n")
	//repl.Start(os.Stdin, os.Stdout)

	data, err := os.ReadFile("src/test_programs/reactivity.flow")
	if err != nil {
		panic(err)
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		panic(p.Errors())
	}
	env := object.NewEnvironment()
	evaluated := eval.Eval(program, env)
	print(evaluated.Inspect())
}
