package main

import (
	"fmt"
	"log"
	"os"

	"Flow/src/eval"
	"Flow/src/lexer"
	"Flow/src/object"
	"Flow/src/parser"
)

func main() {
	filePath := os.Args[1]
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(fmt.Errorf("could not open %s, %w", filePath, err))
	}
	p := parser.New(lexer.New(string(data)))
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		for i := 0; i < len(p.Errors()); i++ {
			log.Print(fmt.Errorf("parsing error: %w", p.Errors()[i]))
		}
		panic(fmt.Errorf("errors while parsing, could not evaluate"))
	}
	env := object.NewEnvironment()
	evaluated := eval.Eval(program, env)
	fmt.Println(evaluated.Inspect())
}
