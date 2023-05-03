package lexer

import (
	"github.com/GzArthur/interpreter/token"
)

// Lexer lexical analysing struct
type Lexer struct {
	input    string // the string to lexer
	ch       byte   // current character
	position int    // current position
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input: input,
	}
	// initialize the input string which preforms lexical parsing
	if len(l.input) > 0 {
		l.ch = l.input[l.position]
	}
	return l
}

// ReadToken read the complete one token from the input string which performs lexical parsing
func (l *Lexer) ReadToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '/':
		tok = token.NewToken(token.SLASH, l.ch)
	case '*':
		tok = token.NewToken(token.ASTERISK, l.ch)
	case '<':
		tok = token.NewToken(token.LT, l.ch)
	case '>':
		tok = token.NewToken(token.GT, l.ch)
	case ';':
		tok = token.NewToken(token.SEMICOLON, l.ch)
	case ',':
		tok = token.NewToken(token.COMMA, l.ch)
	case '{':
		tok = token.NewToken(token.LBRACE, l.ch)
	case '}':
		tok = token.NewToken(token.RBRACE, l.ch)
	case '(':
		tok = token.NewToken(token.LPAREN, l.ch)
	case ')':
		tok = token.NewToken(token.RPAREN, l.ch)
	case '+':
		tok = token.NewToken(token.PLUS, l.ch)
	case '-':
		tok = token.NewToken(token.MINUS, l.ch)
	case '=':
		if l.peekNextCharacter() != '=' {
			tok = token.NewToken(token.ASSIGN, l.ch)
		} else {
			l.readNextCharacter()
			tok = token.NewToken(token.EQ, "==")
		}
	case '!':
		if l.peekNextCharacter() != '=' {
			tok = token.NewToken(token.BANG, l.ch)
		} else {
			l.readNextCharacter()
			tok = token.NewToken(token.NOT_EQ, "!=")
		}
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
		} else {
			tok = token.NewToken(token.ILLEGAL, l.ch)
		}
	}
	l.readNextCharacter()
	return tok
}

// readNextCharacter set currentCharacter value and move relation position
func (l *Lexer) readNextCharacter() {
	l.position++
	if l.position >= len(l.input) {
		// the position out of bounds
		l.ch = 0
	} else {
		l.ch = l.input[l.position]
	}
}

// peekNextCharacter peek next character for some lexeral unit like == or !=
func (l *Lexer) peekNextCharacter() byte {
	nextPos := l.position + 1
	if nextPos >= len(l.input) {
		// the position out of bounds
		return 0
	}
	return l.input[nextPos]
}

// readNumber read the complete one number
func (l *Lexer) readNumber() string {
	startPos := l.position
	for {
		if nc := l.peekNextCharacter(); !isDigit(nc) {
			break
		}
		l.readNextCharacter()
	}
	return l.input[startPos : l.position+1]
}

// readIdentifier read the complete one identifier
func (l *Lexer) readIdentifier() string {
	startPos := l.position
	for {
		if nc := l.peekNextCharacter(); !isLetter(nc) {
			break
		}
		l.readNextCharacter()
	}
	return l.input[startPos : l.position+1]
}

// isLetter judge current character is or isn't letter.
// _ allow to be used in identifier and keyword so _ alse be considered a letter
func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

// isDigit judge current character is or isn't digit
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// skipWhitespace skip whitespace
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readNextCharacter()
	}
}
