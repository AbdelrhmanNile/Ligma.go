package runtime

import (
	"bytes"
	"ligma/token"
	"strings"
)

// ---- Start Identifier Block ----
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) Accept(v ExpressionVisitor) LigmaObject {
	return v.VisitIdentifier(i)
}
func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }
// ---- End Identifier Block ----


// ---- Start IntegerLiteral Block ----
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) Accept(v ExpressionVisitor) LigmaObject {
	return v.VisitIntegerLiteral(il)
}
func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }
// ---- End IntegerLiteral Block ----

// ---- Start FloatLiteral Block ----
type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) Accept(v ExpressionVisitor) LigmaObject {
	return v.VisitFloatLiteral(fl)
}
func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

// ---- Start Boolean Block ----
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) Accept(v ExpressionVisitor) LigmaObject {
	return v.VisitBoolean(b)
}
func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }
// ---- End Boolean Block ----

// ---- Start Null Block ----
type Null struct {
	Token token.Token
}

func (n *Null) Accept(v ExpressionVisitor) LigmaObject {
	return v.VisitNull(n)
}
func (n *Null) expressionNode()      {}
func (n *Null) TokenLiteral() string { return n.Token.Literal }
func (n *Null) String() string       { return n.Token.Literal }

// ---- Start StringLiteral Block ----
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) Accept(v ExpressionVisitor) LigmaObject {
	return v.VisitStringLiteral(sl)
}
func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }
// ---- End StringLiteral Block ----

// ---- Start FunctionLiteral Block ----
type FunctionLiteral struct {
	Token      token.Token // the 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) Accept(v ExpressionVisitor) LigmaObject {
	return v.VisitFunctionLiteral(fl)
}
func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}
// ---- End FunctionLiteral Block ----

// ---- Start ListLiteral Block ----
type ListLiteral struct {
	Token    token.Token // the '[' token
	Elements []Expression
}

func (ll *ListLiteral) Accept(v ExpressionVisitor) LigmaObject {
	return v.VisitListLiteral(ll)
}
func (ll *ListLiteral) expressionNode()      {}
func (ll *ListLiteral) TokenLiteral() string { return ll.Token.Literal }
func (ll *ListLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range ll.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
// ---- End ListLiteral Block ----
