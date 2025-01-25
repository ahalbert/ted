package main

import (
	"fmt"
	"io"
	"os"
	"regexp"

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

	parsedFSA, errors := p.ParseFSA()
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Println(err)
		}
		os.Exit(1)
	}

	if flags.Flags.DebugMode {
		io.WriteString(os.Stdout, parsedFSA.String())
	}

	variables := make(map[string]string)
	if flags.Flags.Seperator == "" {
		variables["$RS"] = "\n"
	} else {
		variables["$RS"] = flags.Flags.Seperator
	}

	if flags.Flags.NoPrint {
		variables["$PRINTMODE"] = "noprint"
	} else {
		variables["$PRINTMODE"] = "print"
	}

	for _, varstring := range flags.Flags.Variables {
		re, err := regexp.Compile("(.*?)=(.*)")
		if err != nil {
			panic("regex compile error")
		}
		matches := re.FindStringSubmatch(varstring)
		if matches != nil {
			variables[matches[1]] = matches[2]
		} else {
			panic("unparsable variable --var " + varstring)
		}
	}

	r := runner.NewRunner(parsedFSA, variables)
	if len(flags.Flags.InputFiles) > 0 {
		for _, infile := range flags.Flags.InputFiles {
			reader, err := os.Open(infile)
			if err != nil {
				panic("Input file " + infile + " not found")
			}
			r.RunFSAFromFile(reader, os.Stdout)
		}
	} else {
		stdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		str := string(stdin)
		if len(str) > 0 {
			str = str[:len(str)-1]
		}
		r.RunFSAFromString(str, os.Stdout)
	}
}
