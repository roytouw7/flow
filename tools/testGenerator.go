package main

import (
	"math/rand"
	"os"
	"time"
)

var validTokens = []string{
	" ",
	"\n",
	"\t",
	"\r",
	"=",
	"==",
	";",
	"(",
	")",
	",",
	"+",
	"-",
	"*",
	"/",
	"<",
	">",
	"!",
	"!=",
	"fn",
	"let",
	"true",
	"false",
	"if",
	"else",
	"return",
}

var validIdentifierRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMOPQRSTUVWXZ_")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randChar() rune {
	return validIdentifierRunes[rand.Intn(len(validIdentifierRunes))]
}

func generateIdentifier(min, max int) string {
	var result = make([]byte, 0)
	identifierLength := rand.Intn(max - min) + min

	for i := 0; i < identifierLength; i++ {
		result = append(result, byte(randChar()))
	}

	return string(result)
}

func getRandomKeyword() string {
	return validTokens[rand.Intn(len(validTokens))]
}

func getKeywordOrIdentifier() string {
	if rand.Intn(4) == 0 {
		return generateIdentifier(1, 15)
	}
	return getRandomKeyword()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var code string
	for i := 0; i < 100000; i ++ {
		code += getKeywordOrIdentifier()
	}

	f, err := os.Create("output.flow")
	defer f.Close()
	checkErr(err)

	_, err = f.WriteString(code)
	checkErr(err)
}
