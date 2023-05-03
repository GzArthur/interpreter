package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/GzArthur/interpreter/lexer"
	"github.com/GzArthur/interpreter/token"
)

const PROMPT = ">> "

// StartREPL Read-Eval-Print Loop
func StartREPL(input io.Reader, output io.Writer) {
	// read the user's input from the input stream
	scanner := bufio.NewScanner(input)
	for {
		fmt.Fprint(output, PROMPT)
		// read a line of user input
		if ok := scanner.Scan(); !ok {
			return
		}
		inputContext := scanner.Text()
		l := lexer.NewLexer(inputContext)

		for tok := l.ReadToken(); tok.Type != token.EOF; tok = l.ReadToken() {
			// %+v will output struct type info
			fmt.Fprintf(output, "%+v\n", tok)
		}
	}
}
