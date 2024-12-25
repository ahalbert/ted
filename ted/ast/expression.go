package ast

import (
	"bytes"
	"strconv"
	"strings"
)

type Expression interface {
	Node
	expressionNode()
}

type Identifier struct {
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string  { return i.Value }

type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string  { return sl.Value }

type Boolean struct {
	Value bool
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) String() string {
	if b.Value {
		return "true"
	}
	return "false"
}

type IntegerLiteral struct {
	Value int
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) String() string  { return strconv.Itoa(il.Value) }

type PrefixExpression struct {
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expressionNode() {}
func (oe *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")

	return out.String()
}

type FunctionLiteral struct {
	Parameters []*Identifier
	Body       Action
}

func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

type CallExpression struct {
	Function  Expression // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}
