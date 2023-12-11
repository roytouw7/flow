package parser

import (
	"fmt"
	"testing"

	"Flow/src/ast"
	"Flow/src/lexer"
)

var tests = []string{
	`let a = 0;`,
	`let add = (x, y) => { return x + y }; add(1, 7);`,
	`let a = 0; let b = 10; let add = (x, y) => { return x + y}; let c = add(a, b); c;`,
	`let a = 0; let b = 1; let c = 2; let d = 3; let e = 4; let f = 5; let g = 6; let h = 7; let i = 8; let j = 9;`,
	`let max = (x, y) => { if (x > y) { return x; } return y; }; max(3, 7);`,
}

var program *ast.Program

func BenchmarkParseProgram(b *testing.B) {
	for testIndex, tt := range tests {
		b.Run(fmt.Sprintf("parse program %d ", testIndex), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				testParser := New(lexer.New(tt))
				p := testParser.ParseProgram()
				if len(testParser.Errors()) > 0 {
					b.Fail()
				}
				program = p
			}
		})
	}
}
