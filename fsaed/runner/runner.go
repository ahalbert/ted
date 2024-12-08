package runner

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/ahalbert/fsaed/fsaed/ast"
	"github.com/ahalbert/fsaed/fsaed/flags"
	"github.com/ahalbert/fsaed/fsaed/parser"
	"github.com/rwtodd/Go.Sed/sed"
)

type Runner struct {
	States                         map[string]*State
	Variables                      map[string]string
	fsa                            ast.FSA
	StartState                     string
	CurrState                      string
	DidTransition                  bool
	DidStartCaptureOnUnderscoreVar bool
	CaptureMode                    string
	CaptureVar                     string
	CurrLine                       string
	parser                         *parser.Parser
}

type State struct {
	StateName string
	Actions   []ast.Action
}

func NewRunner(fsa ast.FSA, p *parser.Parser) *Runner {
	r := &Runner{States: make(map[string]*State), Variables: make(map[string]string)}
	r.parser = p
	r.States["0"] = newState("0")
	r.Variables["$_"] = ""

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
	if flags.Flags.NoPrint {
		r.CaptureMode = "capture"
		r.CaptureVar = "$NULL"
	} else {
		r.CaptureMode = "nocapture"
	}

	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanLines)
	r.CurrState = r.StartState
	for scanner.Scan() {
		r.CurrLine = scanner.Text()
		r.clearAndSetVariable("$@", r.CurrLine+"\n")
		if r.CaptureMode != "underscore" && !(r.CaptureVar == "$_" && r.CaptureMode == "capture") {
			r.clearAndSetVariable("$_", r.CurrLine+"\n")
		}
		r.DidTransition = false
		state, ok := r.States[r.CurrState]
		if !ok {
			panic("missing state:" + r.CurrState)
		}
		for _, action := range state.Actions {
			r.doAction(action)
			if r.DidTransition {
				break
			}
		}

		if r.CaptureMode == "underscore" {
			r.CaptureMode = "capture"
		} else if r.CaptureMode == "capture" {
			r.appendToVariable(r.CaptureVar, r.getVariable("$@"))
		} else if r.CaptureMode == "temp" {
			r.CaptureMode = "nocapture"
		} else if !flags.Flags.NoPrint {
			io.WriteString(os.Stdout, r.getVariable("$_"))
			r.clearAndSetVariable("$_", "")
		} else {
			r.clearAndSetVariable("$_", "")
		}
	}
}

func (r *Runner) getVariable(key string) string {
	val, ok := r.Variables[key]
	if !ok {
		panic("Attempted to reference non-existent variable " + key)
	}
	return val
}

func (r *Runner) appendToVariable(key string, apendee string) string {
	if key == "$NULL" {
		return ""
	}
	val, ok := r.Variables[key]
	if !ok {
		val = ""
	}
	val = val + apendee
	r.Variables[key] = val
	return val
}

func (r *Runner) clearAndSetVariable(key string, toset string) {
	r.Variables[key] = toset
}

func (r *Runner) doTransition(newState string) {
	r.CurrState = newState
	r.DidTransition = true
}

func (r *Runner) applyVariablesToString(input string) string {
	var output bytes.Buffer
	t := template.Must(template.New("").Parse(input))
	t.Execute(&output, r.Variables)
	return output.String()
}

func (r *Runner) doAction(action ast.Action) {
	switch action.(type) {
	case *ast.ActionBlock:
		r.doActionBlock(action.(*ast.ActionBlock))
	case *ast.RegexAction:
		r.doRegexAction(action.(*ast.RegexAction))
	case *ast.DoSedAction:
		r.doSedAction(action.(*ast.DoSedAction))
	case *ast.GotoAction:
		r.doGotoAction(action.(*ast.GotoAction))
	case *ast.PrintAction:
		r.doPrintAction(action.(*ast.PrintAction))
	case *ast.StartStopCaptureAction:
		r.doStartStopCapture(action.(*ast.StartStopCaptureAction))
	case *ast.CaptureAction:
		r.doCaptureAction(action.(*ast.CaptureAction))
	case *ast.ClearAction:
		r.doClearAction(action.(*ast.ClearAction))
	case nil:
		r.doNoOp()
	default:
		panic("Unknown Action!")
	}
}

func (r *Runner) doActionBlock(block *ast.ActionBlock) {
	for _, action := range block.Actions {
		r.doAction(action)
	}
}

func (r *Runner) doRegexAction(action *ast.RegexAction) {
	rule := r.applyVariablesToString(action.Rule)
	re, err := regexp.Compile(rule)
	if err != nil {
		panic("regexp error, supplied: " + action.Rule + "\n formatted as: " + rule)
	}
	if re.MatchString(r.CurrLine) {
		r.doAction(action.Action)
	}
}

func (r *Runner) doSedAction(action *ast.DoSedAction) {
	command := r.applyVariablesToString(action.Command)
	engine, err := sed.New(strings.NewReader(command))
	if err != nil {
		panic("error building sed engine with command: '" + action.Command + "'\n formatted as: '" + command + "'")
	}
	result, err := engine.RunString(r.getVariable(action.Variable))
	if err != nil {
		panic("error running sed")
	}
	r.clearAndSetVariable(action.Variable, result)
}

func (r *Runner) doGotoAction(action *ast.GotoAction) {
	if action.Target == "" {
		r.CurrState = "0"
	} else {
		r.CurrState = action.Target
	}
	r.DidTransition = true
}

func (r *Runner) doPrintAction(action *ast.PrintAction) {
	io.WriteString(os.Stdout, r.getVariable(action.Variable))
}

func (r *Runner) doCaptureAction(action *ast.CaptureAction) {
	r.appendToVariable(action.Variable, r.getVariable("$@"))
	r.CaptureMode = "temp"
}

func (r *Runner) doStartStopCapture(action *ast.StartStopCaptureAction) {
	if action.Command == "start" {
		if action.Variable == "$_" {
			r.clearAndSetVariable("$_", "")
		}
		r.CaptureMode = "capture"
		r.CaptureVar = action.Variable
	} else if action.Command == "stop" {
		r.CaptureMode = "nocapture"
	} else {
		panic("unknown command: " + action.Command + " in start/stop action")
	}
}

func (r *Runner) doClearAction(action *ast.ClearAction) {
	if action.Variable == "$_" {
		r.CaptureMode = "nocapture"
	}
	r.clearAndSetVariable(action.Variable, "")
}

func (r *Runner) doNoOp() {}
