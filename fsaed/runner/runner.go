package runner

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/ahalbert/fsaed/fsaed/ast"
	"github.com/ahalbert/fsaed/fsaed/flags"
	"github.com/ahalbert/fsaed/fsaed/parser"
	"github.com/rwtodd/Go.Sed/sed"
)

type Runner struct {
	States        map[string]*State
	Variables     map[string]string
	fsa           ast.FSA
	StartState    string
	CurrState     string
	didTransition bool
	CurrLine      string
	parser        *parser.Parser
}

type State struct {
	StateName string
	Actions   []ast.Action
}

func NewRunner(fsa ast.FSA, p *parser.Parser) *Runner {
	r := &Runner{States: make(map[string]*State), Variables: make(map[string]string)}
	r.parser = p
	r.States["0"] = newState("0")
	for _, varstring := range flags.Flags.Variables {
		re, err := regexp.Compile("(.*?)=(.*)")
		if err != nil {
			panic("regex compile error")
		}
		matches := re.FindStringSubmatch(varstring)
		if matches != nil {
			r.Variables[matches[1]] = matches[2]
		} else {
			panic("unparsable variable --var " + varstring)
		}
	}
	for _, statement := range fsa.Statements {
		switch statement.(type) {
		case *ast.StateStatement:
			if r.StartState == "" {
				r.StartState = statement.(*ast.StateStatement).StateName
			}
			r.processStateStatement(statement.(*ast.StateStatement))
		}
	}
	return r
}

func (r *Runner) processStateStatement(statement *ast.StateStatement) {
	_, ok := r.States[statement.StateName]
	if !ok {
		r.States[statement.StateName] = newState(statement.StateName)
	}
	state, _ := r.States[statement.StateName]
	state.addRule(statement.Action)

}

func newState(stateName string) *State {
	return &State{StateName: stateName,
		Actions: []ast.Action{},
	}
}

func (s *State) addRule(action ast.Action) {
	s.Actions = append(s.Actions, action)
}

func (r *Runner) RunFSA(input io.Reader) {
	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanLines)
	r.CurrState = r.StartState
	for scanner.Scan() {
		r.CurrLine = scanner.Text()
		r.didTransition = false
		state, ok := r.States[r.CurrState]
		if !ok {
			panic("missing state:" + r.CurrState)
		}
		for _, action := range state.Actions {
			r.doAction(action)
			if r.didTransition {
				break
			}
		}
		if !flags.Flags.NoPrint {
			io.WriteString(os.Stdout, r.CurrLine+"\n")
		}
	}
}

func (r *Runner) doTransition(newState string) {
	r.CurrState = newState
	r.didTransition = true
}

func (r *Runner) doAction(action ast.Action) {
	switch action.(type) {
	case ast.RegexAction:
		r.doRegexAction(action.(ast.RegexAction))
	case ast.DoSedAction:
		r.doSedAction(action.(ast.DoSedAction))
	case ast.GotoAction:
		r.doGotoAction(action.(ast.GotoAction))
	case ast.PrintAction:
		r.doPrintAction(action.(ast.PrintAction))
	case nil:
		r.doNoOp()
	default:
		panic("Unknown Action!")
	}
}

func (r *Runner) doRegexAction(action ast.RegexAction) {
	rule := r.applyVariablesToString(action.Rule)
	re, err := regexp.Compile(rule)
	if err != nil {
		panic("regexp error, supplied: " + action.Rule + "\n formatted as: " + rule)
	}
	if re.MatchString(r.CurrLine) {
		r.doAction(action.Action)
	}
}

func (r *Runner) applyVariablesToString(input string) string {
	var output bytes.Buffer
	t := template.Must(template.New("").Parse(input))
	t.Execute(&output, r.Variables)
	return output.String()
}

func (r *Runner) doSedAction(action ast.DoSedAction) {
	command := r.applyVariablesToString(action.Command)
	engine, err := sed.New(strings.NewReader(command))
	if err != nil {
		panic("error building sed engine with command: " + action.Command + "\n formatted as: " + command)
	}
	r.CurrLine, err = engine.RunString(r.CurrLine)
	r.CurrLine = strings.TrimSuffix(r.CurrLine, "\n")
	if err != nil {
		panic("error running sed")
	}
}

func (r *Runner) doGotoAction(action ast.GotoAction) {
	if action.Target == strconv.Itoa(r.parser.AnonymousStates) {
		r.CurrState = "0"
	} else {
		r.CurrState = action.Target
	}
	r.didTransition = true
}

func (r *Runner) doPrintAction(action ast.PrintAction) {
	io.WriteString(os.Stdout, r.CurrLine+"\n")
}

func (r *Runner) doNoOp() {}
