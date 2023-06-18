package lexer

import "github.com/GzArthur/interpreter/token"

// Lexer lexical analysing struct
type Lexer struct {
	textToBeParsed string // the string to lexer
	currChar       byte   // current character
	currPosition   int    // current currPosition
}

func New(input string) *Lexer {
	l := &Lexer{
		textToBeParsed: input,
	}
	// initialize the textToBeParsed string which preforms lexical parsing
	if len(l.textToBeParsed) > 0 {
		l.currChar = l.textToBeParsed[l.currPosition]
	}
	return l
}

// ReadToken read the complete one token from the textToBeParsed string which performs lexical parsing
func (l *Lexer) ReadToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.currChar {
	case '/':
		tok = token.New(token.SLASH, l.currChar)
	case '*':
		tok = token.New(token.ASTERISK, l.currChar)
	case '<':
		tok = token.New(token.LT, l.currChar)
	case '>':
		tok = token.New(token.GT, l.currChar)
	case ';':
		tok = token.New(token.SEMICOLON, l.currChar)
	case ',':
		tok = token.New(token.COMMA, l.currChar)
	case '{':
		tok = token.New(token.LBRACE, l.currChar)
	case '}':
		tok = token.New(token.RBRACE, l.currChar)
	case '(':
		tok = token.New(token.LPAREN, l.currChar)
	case ')':
		tok = token.New(token.RPAREN, l.currChar)
	case '+':
		tok = token.New(token.PLUS, l.currChar)
	case '-':
		tok = token.New(token.MINUS, l.currChar)
	case '=':
		if l.peekNextCharacter() == '=' {
			l.readNextCharacter()
			tok = token.New(token.EQ, "==")
		} else {
			tok = token.New(token.ASSIGN, l.currChar)
		}
	case '!':
		if l.peekNextCharacter() == '=' {
			l.readNextCharacter()
			tok = token.New(token.NOT_EQ, "!=")
		} else {
			tok = token.New(token.BANG, l.currChar)
		}
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.currChar) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
		} else if isDigit(l.currChar) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
		} else {
			tok = token.New(token.ILLEGAL, l.currChar)
		}
	}
	l.readNextCharacter()
	return tok
}

// skipWhitespace skip whitespace
func (l *Lexer) skipWhitespace() {
	for l.currChar == ' ' || l.currChar == '\t' || l.currChar == '\n' || l.currChar == '\r' {
		l.readNextCharacter()
	}
}

// readNumber read the complete one number
func (l *Lexer) readNumber() string {
	startPos := l.currPosition
	for {
		if nc := l.peekNextCharacter(); !isDigit(nc) {
			break
		}
		l.readNextCharacter()
	}
	return l.textToBeParsed[startPos : l.currPosition+1]
}

// readIdentifier read the complete one identifier
func (l *Lexer) readIdentifier() string {
	startPos := l.currPosition
	for {
		if nc := l.peekNextCharacter(); !isLetter(nc) {
			break
		}
		l.readNextCharacter()
	}
	return l.textToBeParsed[startPos : l.currPosition+1]
}

// readNextCharacter set currentCharacter value and move relation currPosition
func (l *Lexer) readNextCharacter() {
	l.currPosition++
	if l.currPosition >= len(l.textToBeParsed) {
		// the currPosition out of bounds
		l.currChar = 0
	} else {
		l.currChar = l.textToBeParsed[l.currPosition]
	}
}

// peekNextCharacter peek next character for some lexeral unit like == or !=
func (l *Lexer) peekNextCharacter() byte {
	nextPos := l.currPosition + 1
	if nextPos >= len(l.textToBeParsed) {
		// the currPosition out of bounds
		return 0
	}
	return l.textToBeParsed[nextPos]
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
