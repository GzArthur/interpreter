package token

const (
	ILLEGAL = "IILLEGAL" // unknown
	EOF     = "EOF"      // end of file

	// Identifier
	IDENTIFIER = "IDENTIFIER" // variable name: x、y, function name: max、add

	// Literals
	INT = "INT"

	// Operator
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"
	EQ       = "=="
	NOT_EQ   = "!="

	// Delimiter
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keyword
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

// keywords map
var keywords = map[string]Type{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// Type lexical unit type
type Type string

// Token lexical unit token
type Token struct {
	Type    Type
	Literal string // lexical unit literals
}

func New[T string | byte](tokenType Type, ch T) Token {
	return Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

// LookupIdent lookup ident from keyword map
func LookupIdent(key string) Type {
	if tok, ok := keywords[key]; ok {
		return tok
	}
	return IDENTIFIER
}
