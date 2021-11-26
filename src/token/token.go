package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	//	Identifiers and literals
	IDENT = "IDENT"
	INT   = "INT"

	//	Operators
	ASSIGN = "="
	PLUS   = "+"
	MINUS = "-"
	BANG = "!"
	ASTERISK = "*"
	SLASH = "/"

	EQ = "=="
	NOT_EQ = "!="

	LT = "<"
	GT = ">"

	//	Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	NEWLINE = "\n"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	//	Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE = "TRUE"
	FALSE = "FALSE"
	IF = "IF"
	ELSE = "ELSE"
	RETURN = "RETURN"
)

type Type string

type Token struct {
	Type    Type
	Literal string
}

var keywords = map[string]Type {
	"fn": FUNCTION,
	"let": LET,
	"true": TRUE,
	"false": FALSE,
	"if": IF,
	"else": ELSE,
	"return": RETURN,
}

// LookupIdentType checks whether input is reserved keyword or identifier
func LookupIdentType(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}

func New(t Type, ch byte) Token {
	return Token{
		Type:    t,
		Literal: string(ch),
	}
}

func IsLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func IsDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}