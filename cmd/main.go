package main

import (
	"os"

	"github.com/ahalbert/fsaed/fsaed/flags"
	"github.com/ahalbert/fsaed/fsaed/lexer"
	"github.com/ahalbert/fsaed/fsaed/parser"
	"github.com/ahalbert/fsaed/fsaed/runner"
	"github.com/alexflint/go-arg"
)

func main() {

	arg.MustParse(&flags.Flags)
	// l := lexer.New(flags.Flags.Program)
	// for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
	// 	fmt.Printf("%+v\n", tok)
	// }

	l := lexer.New(flags.Flags.Program)
	p := parser.New(l)

	parsedFSA := p.ParseFSA()

	// io.WriteString(os.Stdout, parsedFSA.String())

	r := runner.NewRunner(parsedFSA, p)
	r.RunFSA(os.Stdin)

}
