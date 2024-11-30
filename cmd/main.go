package main

import (
	"os"

	"github.com/ahalbert/fsaed/fsaed/lexer"
	"github.com/ahalbert/fsaed/fsaed/parser"
	"github.com/ahalbert/fsaed/fsaed/runner"
)

func main() {
	program := os.Args[1]
	// l := lexer.New(program)
	// for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
	// 	fmt.Printf("%+v\n", tok)
	// }

	l := lexer.New(program)
	p := parser.New(l)

	parsedFSA := p.ParseFSA()

	// io.WriteString(os.Stdout, parsedFSA.String())

	r := runner.NewRunner(parsedFSA)
	r.RunFSA(os.Stdin)

}
