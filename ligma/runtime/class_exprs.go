package runtime

import (
	"bytes"
	"ligma/token"
)

// ---- Start Class Block ----
type Class struct {
	Token token.Token
	Name *Identifier
	Superclass *Identifier
	Methods []*DefStatement
}

func (c *Class) Accept(v StatementVisitor) LigmaObject {
	return v.VisitClassStatement(c)
}
func (c *Class) statementNode()       {}
func (c *Class) TokenLiteral() string { return c.Token.Literal }
func (c *Class) String() string {
	var out bytes.Buffer

	out.WriteString("class ")
	out.WriteString(c.Name.String())
	out.WriteString(" {\n")

	for _, m := range c.Methods {
		out.WriteString(m.String())
	}

	out.WriteString("\n}")
	return out.String()
}
// ---- End Class Block ----

// ---- Start GetExpression Block ----
type GetExpression struct {
	Token token.Token
	Object Expression
	Property *Identifier
}

func (ge *GetExpression) Accept(v ExpressionVisitor) LigmaObject {
	return v.VisitGetExpression(ge)
}
func (ge *GetExpression) expressionNode()      {}
func (ge *GetExpression) TokenLiteral() string { return ge.Token.Literal }
func (ge *GetExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ge.Object.String())
	out.WriteString(".")
	out.WriteString(ge.Property.String())

	return out.String()
}
// ---- End GetExpression Block ----

// ---- Start SetExpression Block ----
type SetExpression struct {
	Token token.Token
	Object Expression
	Property *Identifier
	Value Expression
}

func (se *SetExpression) Accept(v ExpressionVisitor) LigmaObject {
	return v.VisitSetExpression(se)
}
func (se *SetExpression) expressionNode()      {}
func (se *SetExpression) TokenLiteral() string { return se.Token.Literal }
func (se *SetExpression) String() string {
	var out bytes.Buffer

	out.WriteString(se.Object.String())
	out.WriteString(".")
	out.WriteString(se.Property.String())
	out.WriteString(" = ")
	out.WriteString(se.Value.String())

	return out.String()
}

// ---- Start Self Block ----
type Self struct {
	Token token.Token
}

func (s *Self) Accept(v ExpressionVisitor) LigmaObject {
	return v.VisitSelfExpression(s)
}
func (s *Self) expressionNode()      {}
func (s *Self) TokenLiteral() string { return s.Token.Literal }
func (s *Self) String() string       { return s.Token.Literal }
// ---- End Self Block ----

// ---- Start Super Block ----
type Super struct {
	Token token.Token
	Method *Identifier
}

func (s *Super) Accept(v ExpressionVisitor) LigmaObject {
	return v.VisitSuper(s)
}
func (s *Super) expressionNode()      {}
func (s *Super) TokenLiteral() string { return s.Token.Literal }
func (s *Super) String() string       { return s.Token.Literal }