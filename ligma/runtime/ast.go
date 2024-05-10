package runtime

import (
	"bytes"
)

// NodeType
const (
	PROGRAM = "Program"
	Def = "DefStatement"
	IDENT = "Identifier"
)


// Visitor interface
type ExpressionVisitor interface {
	VisitPrefixExpression(*PrefixExpression) LigmaObject
	VisitInfixExpression(*InfixExpression) LigmaObject
	VisitIfExpression(*IfExpression) LigmaObject
	VisitCallExpression(*CallExpression) LigmaObject
	VisitIndexExpression(*IndexExpression) LigmaObject
	VisitAssignExpression(*AssignExpression) LigmaObject
	VisitIdentifier(*Identifier) LigmaObject
	VisitIntegerLiteral(*IntegerLiteral) LigmaObject
	VisitFloatLiteral(*FloatLiteral) LigmaObject
	VisitBoolean(*Boolean) LigmaObject
	VisitNull(*Null) LigmaObject
	VisitListLiteral(*ListLiteral) LigmaObject
	VisitStringLiteral(*StringLiteral) LigmaObject
	VisitFunctionLiteral(*FunctionLiteral) LigmaObject
	VisitMapLiteral(*MapLiteral) LigmaObject
	VisitGetExpression(*GetExpression) LigmaObject
	VisitSetExpression(*SetExpression) LigmaObject
	VisitSelfExpression(*Self) LigmaObject
	VisitSuper(*Super) LigmaObject
}

type StatementVisitor interface {
	VisitDefStatement(*DefStatement) LigmaObject
	VisitReturnStatement(*ReturnStatement) LigmaObject
	VisitExpressionStatement(*ExpressionStatement) LigmaObject
	VisitBlockStatement(*BlockStatement)  LigmaObject
	VisitClassStatement(*Class) LigmaObject
	VisitWhileStatement(*WhileStatement) LigmaObject
}


type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
	Accept(StatementVisitor) LigmaObject
}

type Expression interface {
	Node
	expressionNode()
	Accept(ExpressionVisitor) LigmaObject
}

// ---- Start Program Block ----
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// ---- End Program Block ----






