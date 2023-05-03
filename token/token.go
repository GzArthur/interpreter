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

// TokenType lexical unit type
type TokenType string

// Token lexical unit token
type Token struct {
	Type    TokenType
	Literal string // lexical unit literals
}

func NewToken[T string | byte](tokenType TokenType, ch T) Token {
	return Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

// keywords map
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// LookupIdent lookup ident from keyword map
func LookupIdent(key string) TokenType {
	if tok, ok := keywords[key]; ok {
		return tok
	}
	return IDENTIFIER
}
