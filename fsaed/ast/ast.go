package ast

import (
	"bytes"
)

// The base Node interface
type Node interface {
	String() string
}

// All statement nodes implement this
type Statement interface {
	Node
	statementNode()
}

// All expression nodes implement this
type Action interface {
	Node
}

type FSA struct {
	Statements []Statement
}

func (fsa *FSA) String() string {
	var out bytes.Buffer

	for _, s := range fsa.Statements {
		out.WriteString(s.String() + "\n")
	}

	return out.String()
}

type StateStatement struct {
	StateName string
	Action    Action
}

func (ss *StateStatement) statementNode() {}
func (ss *StateStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ss.StateName + ":" + ss.Action.String())
	return out.String()
}

type ActionBlock struct {
	Actions []Action
}

func (ab *ActionBlock) String() string {
	var out bytes.Buffer
	out.WriteString("{ ")
	for _, action := range ab.Actions {
		out.WriteString(action.String() + "; ")
	}
	out.WriteString(" }")
	return out.String()
}

type RegexAction struct {
	Rule   string
	Action Action
}

func (ra *RegexAction) String() string {
	var out bytes.Buffer
	out.WriteString("/" + ra.Rule + "/ " + " :: " + (ra.Action).String())
	return out.String()
}

type GotoAction struct {
	Target string
}

func (ga *GotoAction) String() string {
	var out bytes.Buffer
	out.WriteString("goto: " + ga.Target)
	return out.String()
}

type DoSedAction struct {
	Variable string
	Command  string
}

func (da *DoSedAction) getNextAction() Action { return nil }
func (da *DoSedAction) String() string {
	var out bytes.Buffer
	out.WriteString("sed '" + da.Command + "' using var '" + da.Variable + "'")
	return out.String()
}

type PrintAction struct {
	Variable string
}

func (pa *PrintAction) String() string {
	var out bytes.Buffer
	out.WriteString("print '" + pa.Variable + "'")
	return out.String()
}

type PrintLnAction struct {
	Variable string
}

func (pa *PrintLnAction) String() string {
	var out bytes.Buffer
	out.WriteString("print '" + pa.Variable + "'")
	return out.String()
}

type StartStopCaptureAction struct {
	Command  string
	Variable string
}

func (sscp *StartStopCaptureAction) String() string {
	var out bytes.Buffer
	out.WriteString(sscp.Command + " capture into:" + sscp.Variable)
	return out.String()
}

type CaptureAction struct {
	Variable string
}

func (ca *CaptureAction) String() string {
	var out bytes.Buffer
	out.WriteString("temp capture into:" + ca.Variable)
	return out.String()
}

type ClearAction struct {
	Variable string
}

func (ca *ClearAction) String() string {
	var out bytes.Buffer
	out.WriteString("clear:" + ca.Variable)
	return out.String()
}

type AssignAction struct {
	Target       string
	IsIdentifier bool
	Value        string
}

func (aa *AssignAction) String() string {
	var out bytes.Buffer
	out.WriteString("set:" + aa.Target + "'= ")
	if aa.IsIdentifier {
		out.WriteString(aa.Value)
	} else {
		out.WriteString("'" + aa.Value + "'")
	}
	return out.String()
}
