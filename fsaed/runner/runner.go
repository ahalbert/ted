package runner

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/ahalbert/fsaed/fsaed/ast"
	"github.com/rwtodd/Go.Sed/sed"
)

type Runner struct {
	States        map[string]*State
	fsa           ast.FSA
	StartState    string
	CurrState     string
	didTransition bool
	CurrLine      string
}

type State struct {
	StateName   string
	Actions     []ast.Action
	Transitions []*ast.TransitionAction
}

func NewRunner(fsa ast.FSA) *Runner {
	r := &Runner{States: make(map[string]*State)}
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

func RunFSA() {
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
		Actions:     []ast.Action{},
		Transitions: []*ast.TransitionAction{},
	}
}

func (s *State) addRule(action ast.Action) {
	switch action.(type) {
	case *ast.TransitionAction:
		s.Transitions = append(s.Transitions, action.(*ast.TransitionAction))
	default:
		s.Actions = append(s.Actions, action)
	}
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
			panic("missing state")
		}
		for _, action := range state.Actions {
			r.doAction(action)
		}
		for _, rule := range state.Transitions {
			if r.didTransition {
				break
			}
			re, err := regexp.Compile(rule.Rule)
			if err != nil {
				panic("re compile error")
			}
			if re.MatchString(r.CurrLine) {
				r.doTransition(rule.Target)
				if rule.Action != nil {
					r.doAction(rule.Action)
				}
			}
		}
		io.WriteString(os.Stdout, r.CurrLine+"\n")
	}
}

func (r *Runner) doTransition(newState string) {
	r.CurrState = newState
	r.didTransition = true
}

func (r *Runner) doAction(action ast.Action) {
	switch action.(type) {
	case *ast.SedAction:
		r.CurrLine = r.doSedAction(action.(*ast.SedAction), r.CurrLine)
	default:
		panic("Unknown Action!")
	}
}

func (r *Runner) doSedAction(action *ast.SedAction, input string) string {
	engine, err := sed.New(strings.NewReader(action.Command))
	if err != nil {
		panic("error building sed engine")
	}
	output, err := engine.RunString(input)
	if err != nil {
		panic("error running sed")
	}
	return output
}
