package runtime

import (
	"bytes"
	"ligma/token"
)

// ---- Start DefStatement Block ----
type DefStatement struct {
	Token token.Token // the token.DEF token
	Name  *Identifier
	Value Expression
}

func (ls *DefStatement) Accept(v StatementVisitor) LigmaObject {
	return v.VisitDefStatement(ls)
}
func (ls *DefStatement) statementNode()       {}
func (ls *DefStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *DefStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

// ---- End DefStatement Block ----

type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) Accept(v StatementVisitor) LigmaObject {
	return v.VisitReturnStatement(rs)
}
func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")
	return out.String()
}

// ---- Start ExpressionStatement Block ----
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) Accept(v StatementVisitor) LigmaObject {
	return v.VisitExpressionStatement(es)
}
func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}
// ---- End ExpressionStatement Block ----

// ---- Start BlockStatement Block ----
type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) Accept(v StatementVisitor) LigmaObject {
	return v.VisitBlockStatement(bs)
}
func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}
// ---- End BlockStatement Block ----


// ---- Start WhileStatement Block ----
type WhileStatement struct {
	Token token.Token
	Condition Expression
	Body *BlockStatement
}

func (ws *WhileStatement) Accept(v StatementVisitor) LigmaObject {
	return v.VisitWhileStatement(ws)
}
func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Literal }
func (ws *WhileStatement) String() string {
	var out bytes.Buffer

	out.WriteString("while ")
	out.WriteString(ws.Condition.String())
	out.WriteString(" ")
	out.WriteString(ws.Body.String())

	return out.String()
}
// ---- End WhileStatement Block ----
