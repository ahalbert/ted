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
	getNextAction() Action
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

func (ab *ActionBlock) getNextAction() Action { return ab.Actions[0] }
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

func (ra *RegexAction) getNextAction() Action { return ra.Action }
func (ra *RegexAction) String() string {
	var out bytes.Buffer
	out.WriteString("/" + ra.Rule + "/ " + " :: " + (ra.Action).String())
	return out.String()
}

type GotoAction struct {
	Target string
}

func (ga *GotoAction) getNextAction() Action { return nil }
func (ga *GotoAction) String() string {
	var out bytes.Buffer
	out.WriteString("goto: " + ga.Target)
	return out.String()
}

type DoSedAction struct {
	Command string
}

func (da *DoSedAction) getNextAction() Action { return nil }
func (da *DoSedAction) String() string {
	var out bytes.Buffer
	out.WriteString("sed " + da.Command)
	return out.String()
}

type PrintAction struct {
}

func (pa *PrintAction) getNextAction() Action { return nil }
func (pa *PrintAction) String() string {
	var out bytes.Buffer
	out.WriteString("print")
	return out.String()
}
