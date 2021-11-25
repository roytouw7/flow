package lexer

import (
	"Flow/src/token"
	"os"
	"testing"
)

func benchmarkNextToken(input string, b *testing.B) {
	l := New(input)
	i := 0

	for n := 0; n < b.N; n++ {
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			_ = tok.Type
			i++
		}
	}

}

func BenchmarkNextToken1(b *testing.B) {
	data, err := os.ReadFile("test_program.flow")
	if err != nil {
		panic(err)
	}
	benchmarkNextToken(string(data), b)
}

func BenchmarkNextToken10000(b *testing.B) {
	data, err := os.ReadFile("10000.flow")
	if err != nil {
		panic(err)
	}
	benchmarkNextToken(string(data), b)
}