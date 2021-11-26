package repl

import (
	"Flow/src/lexer"
	"Flow/src/token"
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		ok, err := regexp.Match(".*\\.flow", []byte(scanner.Text()))
		checkErr(err)

		var line string

		if ok {
			data, err := os.ReadFile("src/main/" + scanner.Text())
			checkErr(err)
			line = string(data)
		} else {
			line = scanner.Text()
		}

		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
