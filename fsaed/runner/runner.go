package runner

import (
	"bytes"
	"errors"
	"fmt"
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
	States                map[string]*State
	Variables             map[string]string
	StartState            string
	CurrState             string
	DidTransition         bool
	DidResetUnderscoreVar bool
	CaptureMode           string
	CaptureVar            string
	CurrLine              int
	MaxLine               int
	Tape                  io.ReadSeeker
	TapeOffsets           map[int]int64
	ShouldHalt            bool
}

type State struct {
	StateName string
	Actions   []ast.Action
}

func NewRunner(fsa ast.FSA, p *parser.Parser) *Runner {
	r := &Runner{
		States:      make(map[string]*State),
		Variables:   make(map[string]string),
		TapeOffsets: make(map[int]int64),
	}
	r.States["0"] = newState("0")
	r.Variables["$_"] = ""
	r.CurrLine = -1
	r.MaxLine = 0
	r.TapeOffsets[0] = 0

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

func (r *Runner) RunFSA(input io.ReadSeeker) {
	r.Tape = input
	if flags.Flags.NoPrint {
		r.CaptureMode = "capture"
		r.CaptureVar = "$NULL"
	} else {
		r.CaptureMode = "nocapture"
	}

	r.CurrState = r.StartState
	for !r.ShouldHalt {
		r.CurrLine++
		line, err := r.getLine(r.CurrLine)
		if err != nil && errors.Is(err, io.EOF) {
			r.ShouldHalt = true
			break
		} else if err != nil {
			panic(err)
		}
		r.clearAndSetVariable("$@", line)

		if !(r.CaptureVar == "$_" && r.CaptureMode == "capture") {
			r.clearAndSetVariable("$_", r.getVariable("$@"))
			r.DidResetUnderscoreVar = true
		} else {
			r.DidResetUnderscoreVar = false
		}

		r.DidTransition = false

		state, ok := r.States[r.CurrState]
		if !ok {
			panic("missing state:" + r.CurrState)
		}

		if r.ShouldHalt {
			break
		}

		for _, action := range state.Actions {
			if r.DidTransition || r.ShouldHalt {
				break
			}
			r.doAction(action)
		}

		if r.ShouldHalt {
			break
		}

		if r.CaptureMode == "capture" {
			r.appendToVariable(r.CaptureVar, r.getVariable("$@")+"\n")
		} else if r.CaptureMode == "temp" {
			r.CaptureMode = "nocapture"
		} else if !flags.Flags.NoPrint {
			io.WriteString(os.Stdout, r.getVariable("$_")+"\n")
			r.clearAndSetVariable("$_", "")
		} else {
			r.clearAndSetVariable("$_", "")
		}
	}
}

func (r *Runner) getLine(line int) (string, error) {
	if line < 0 {
		return "", fmt.Errorf("line %d less than 0", line)
	}
	offset, ok := r.TapeOffsets[line]
	if !ok {
		linenum := r.MaxLine
		offset, _ := r.TapeOffsets[linenum]
		r.Tape.Seek(offset, 0)
		for linenum < line {
			_, err := r.readToNewline()
			linenum++
			if err != nil {
				return "", err
			}
			offset, err = r.Tape.Seek(0, io.SeekCurrent)
			if err != nil {
				panic(err)
			}
			r.TapeOffsets[linenum] = offset
			r.MaxLine = linenum
		}
	} else {
		r.Tape.Seek(offset, 0)
	}
	return r.readToNewline()
}

func (r *Runner) readToNewline() (string, error) {
	line := ""
	bit := make([]byte, 1)
	for ok := true; ok; ok = (bit[0] != '\n') {
		line += string(bit[0])
		_, err := r.Tape.Read(bit)
		if err != nil {
			return "", err
		}
	}
	return line[1:], nil
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
	case *ast.PrintLnAction:
		r.doPrintLnAction(action.(*ast.PrintLnAction))
	case *ast.StartStopCaptureAction:
		r.doStartStopCapture(action.(*ast.StartStopCaptureAction))
	case *ast.CaptureAction:
		r.doCaptureAction(action.(*ast.CaptureAction))
	case *ast.ClearAction:
		r.doClearAction(action.(*ast.ClearAction))
	case *ast.AssignAction:
		r.doAssignAction(action.(*ast.AssignAction))
	case *ast.MoveHeadAction:
		r.doMoveHeadAction(action.(*ast.MoveHeadAction))
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

	matches := re.FindStringSubmatch(r.getVariable("$@"))
	if matches != nil {
		for idx, match := range matches {
			stridx := "$" + strconv.Itoa(idx)
			r.clearAndSetVariable(stridx, match)
		}
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
	if action.Variable == "$_" && r.CaptureMode != "capture" {
		r.clearAndSetVariable(action.Variable, result[:len(result)-1])
	} else {
		r.clearAndSetVariable(action.Variable, result)
	}
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

func (r *Runner) doPrintLnAction(action *ast.PrintLnAction) {
	io.WriteString(os.Stdout, r.getVariable(action.Variable)+"\n")
}

func (r *Runner) doCaptureAction(action *ast.CaptureAction) {
	if action.Variable == "$_" && r.DidResetUnderscoreVar {
		r.clearAndSetVariable("$_", "")
	}
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

func (r *Runner) doAssignAction(action *ast.AssignAction) {
	if action.IsIdentifier {
		r.Variables[action.Target] = r.getVariable(action.Value)
	} else {
		r.Variables[action.Target] = r.applyVariablesToString(action.Value)
	}
}

func (r *Runner) doMoveHeadAction(action *ast.MoveHeadAction) {
	if action.Command == "fastforward" {
		r.doFastForward(action.Regex)
	} else if action.Command == "rewind" {
		r.doRewind(action.Regex)
	}
}

func (r *Runner) doFastForward(target string) {
	rule := r.applyVariablesToString(target)
	re, err := regexp.Compile(rule)
	if err != nil {
		panic(err)
	}
	line := ""
	for ok := true; ok; ok = (!re.MatchString(line)) {
		r.CurrLine++
		line, err = r.getLine(r.CurrLine)
		if err != nil && errors.Is(err, io.EOF) {
			r.ShouldHalt = true
			break
		} else if err != nil {
			panic(err)
		}
	}
	r.CurrLine--
}

func (r *Runner) doRewind(target string) {
	rule := r.applyVariablesToString(target)
	re, err := regexp.Compile(rule)
	if err != nil {
		panic(err)
	}
	line := ""
	for ok := true; ok; ok = (!re.MatchString(line) && r.CurrLine >= 0) {
		r.CurrLine--
		line, err = r.getLine(r.CurrLine)
		if err != nil {
			panic(err)
		}
	}
	r.CurrLine--
}

func (r *Runner) doNoOp() {}
