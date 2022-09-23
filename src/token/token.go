package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	UNKNOWN = "???"

	//	Identifiers and literals
	IDENT = "IDENT"
	INT   = "INT"

	//	Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	QUESTION = "?"
	COLON    = ":"

	EQ     = "=="
	NOT_EQ = "!="

	LT = "<"
	GT = ">"

	//	Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	NEWLINE   = "\n"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	//	Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

type Type string

type Token struct {
	Type      Type
	Literal   string
	Pos, Line int
}

var keywords = map[string]Type{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// LookupIdentType checks whether input is reserved keyword or identifier
func LookupIdentType(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}

func New(t Type, ch string, pos, line int) *Token {
	return &Token{
		Type:    t,
		Literal: ch,
		Pos:     pos,
		Line:    line,
	}
}

func NewSymbol(t Type, pos, line int) *Token {
	return New(t, string(t), pos, line) // type conversion actually required
}
