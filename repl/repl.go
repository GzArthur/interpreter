package repl

import (
	"bufio"
	"fmt"
	"github.com/GzArthur/interpreter/evaluator"
	"github.com/GzArthur/interpreter/object"
	"github.com/GzArthur/interpreter/parser"
	"io"

	"github.com/GzArthur/interpreter/lexer"
)

const PROMPT = ">> "

// StartREPL Read-Eval-Print Loop
func StartREPL(input io.Reader, output io.Writer) {
	// read the user's input from the input stream
	scanner := bufio.NewScanner(input)
	env := object.NewEnv()
	for {
		// outputs the identity >>  before user input
		if _, err := fmt.Fprint(output, PROMPT); err != nil {
			panic("print >> Error")
		}
		// read a line of user input
		if ok := scanner.Scan(); !ok {
			panic("scan Error")
		}
		inputContext := scanner.Text()
		l := lexer.New(inputContext)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(output, p.Errors())
			continue
		}

		obj := evaluator.Eval(program, env)
		if obj != nil {
			if _, err := io.WriteString(output, fmt.Sprintf("%s\n", obj.Inspect())); err != nil {
				panic("output Error")
			}
		}
	}
}

func printParserErrors(output io.Writer, errors []string) {
	for _, msg := range errors {
		if _, err := io.WriteString(output, fmt.Sprintf("Woops! a parser error has occurred:\n \t%s\n", msg)); err != nil {
			panic("print parser errors Error")
		}
	}
}
