package runtime

import (
	"bytes"
	"ligma/token"
	"strings"
)

// ---- Start PrefixExpression Block ----
type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pr *PrefixExpression) Accept(v ExpressionVisitor) LigmaObject {
	return v.VisitPrefixExpression(pr)
}
func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}
// ---- End PrefixExpression Block ----

// ---- Start InfixExpression Block ----
type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) Accept(v ExpressionVisitor) LigmaObject {
	return v.VisitInfixExpression(ie)
}
func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}
// ---- End InfixExpression Block ----

// ---- Start IfExpression Block ----
type IfExpression struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) Accept(v ExpressionVisitor) LigmaObject {
	return v.VisitIfExpression(ie)
}
func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}
// ---- End IfExpression Block ----

// ---- Start IndexExpression Block ----
type IndexExpression struct {
	Token token.Token // The '[' token
	Left  Expression // object to index
	Index Expression // The index expression
}

func (ie *IndexExpression) Accept(v ExpressionVisitor) LigmaObject {
	return v.VisitIndexExpression(ie)
}
func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

// ---- Start CallExpression Block ----
type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) Accept(v ExpressionVisitor) LigmaObject {
	return v.VisitCallExpression(ce)
}
func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
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
// ---- End CallExpression Block ----

// ---- Start AssignExpression Block ----
type AssignExpression struct {
	Token token.Token
	Name *Identifier
	Value Expression
}

func (ae *AssignExpression) Accept(v ExpressionVisitor) LigmaObject {
	return v.VisitAssignExpression(ae)
}
func (ae *AssignExpression) expressionNode()      {}
func (ae *AssignExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AssignExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ae.Name.String())
	out.WriteString(" = ")
	out.WriteString(ae.Value.String())

	return out.String()
}
// ---- End AssignExpression Block ----
