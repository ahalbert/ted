package main

import (
	"fmt"
	"io"
	"os"

	"github.com/ahalbert/ted/ted/flags"
	"github.com/ahalbert/ted/ted/lexer"
	"github.com/ahalbert/ted/ted/parser"
	"github.com/ahalbert/ted/ted/runner"
	"github.com/ahalbert/ted/ted/token"
	"github.com/alexflint/go-arg"
)

func main() {

	arg.MustParse(&flags.Flags)
	var program string
	if flags.Flags.ProgramFile != "" {
		buf, err := os.ReadFile(flags.Flags.ProgramFile)
		if err != nil {
			panic("FSA File " + flags.Flags.ProgramFile + " not found")
		}
		program = string(buf)
		if flags.Flags.Program != "" {
			flags.Flags.InputFiles = append([]string{flags.Flags.Program}, flags.Flags.InputFiles...)
		}
	} else {
		program = flags.Flags.Program
	}

	if program == "" {
		panic("no FSA supplied")
	}

	if flags.Flags.DebugMode {
		l := lexer.New(program)
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}

	l := lexer.New(program)
	p := parser.New(l)

	parsedFSA := p.ParseFSA()

	if flags.Flags.DebugMode {
		io.WriteString(os.Stdout, parsedFSA.String())
	}

	r := runner.NewRunner(parsedFSA, p)
	if len(flags.Flags.InputFiles) > 0 {
		for _, infile := range flags.Flags.InputFiles {
			reader, err := os.Open(infile)
			if err != nil {
				panic("Input file " + infile + " not found")
			}
			r.RunFSA(reader)
		}
	} else {
		r.RunFSA(os.Stdin)
	}

}
