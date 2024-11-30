package ast

import (
	"bytes"

	"github.com/ahalbert/fsaed/fsaed/token"
)

// The base Node interface
type Node interface {
	TokenLiteral() string
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
	actionNode()
}

type FSA struct {
	Statements []Statement
}

func (fsa *FSA) TokenLiteral() string {
	if len(fsa.Statements) > 0 {
		return fsa.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (fsa *FSA) String() string {
	var out bytes.Buffer

	for _, s := range fsa.Statements {
		out.WriteString(s.String() + "\n")
	}

	return out.String()
}

type StateStatement struct {
	Token     token.Token // The first token in the statment
	StateName string
	Action    Action
}

func (ss *StateStatement) statementNode()       {}
func (ss *StateStatement) TokenLiteral() string { return ss.Token.Literal }
func (ss *StateStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ss.StateName + ":" + ss.Action.String())
	return out.String()
}

type TransitionAction struct {
	Token  token.Token // The first token in the Action
	Rule   string
	Target string
	Action Action
}

func (ta TransitionAction) actionNode()          {}
func (ta TransitionAction) TokenLiteral() string { return ta.Token.Literal }
func (ta TransitionAction) String() string {
	var out bytes.Buffer
	out.WriteString("/" + ta.Rule + "/ " + "-> " + ta.Target)
	if ta.Action != nil {
		out.WriteString(", " + ta.Action.String())
	}
	return out.String()
}

type SedAction struct {
	Token   token.Token // The first token in the Action
	Command string
}

func (sa SedAction) actionNode()          {}
func (sa SedAction) TokenLiteral() string { return sa.Token.Literal }
func (sa SedAction) String() string       { return sa.Command }
