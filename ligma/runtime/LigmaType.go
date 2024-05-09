package runtime

import (
	"bytes"
	"fmt"
	"strings"
)

// LigmaInteger
type LigmaInteger struct {
	Value int64
}

func (i *LigmaInteger) Inspect() string { return fmt.Sprintf("%d", i.Value) }
func (i *LigmaInteger) Type() ObjectType { return INTEGER_OBJ }

// LigmaFloat
type LigmaFloat struct {
	Value float64
}

func (f *LigmaFloat) Inspect() string { return fmt.Sprintf("%f", f.Value) }
func (f *LigmaFloat) Type() ObjectType { return FLOAT_OBJ }


// LigmaBoolean
type LigmaBoolean struct {
	Value bool
}

func (b *LigmaBoolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }
func (b *LigmaBoolean) Type() ObjectType { return BOOLEAN_OBJ }

// LigmaNull
type LigmaNull struct{}

func (n *LigmaNull) Inspect() string { return "null" }
func (n *LigmaNull) Type() ObjectType { return NULL_OBJ }

// LigmaString
type LigmaString struct {
	Value string
}

func (s *LigmaString) Inspect() string { return s.Value }
func (s *LigmaString) Type() ObjectType { return STRING_OBJ }

// LigmaList
type LigmaList struct {
	Elements []LigmaObject
}

func (l *LigmaList) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range l.Elements {
		elements = append(elements, el.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
func (l *LigmaList) Type() ObjectType { return LIST_OBJ }
