package runtime

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"
)

// LigmaInteger
type LigmaInteger struct {
	Value int64
}

func (i *LigmaInteger) Inspect() string { return fmt.Sprintf("%d", i.Value) }
func (i *LigmaInteger) Type() ObjectType { return INTEGER_OBJ }
func (i *LigmaInteger) MapKey() MapKey {
	var key MapKey
	key.Value = uint64(i.Value)
	return MapKey{Type: i.Type(), Value: key.Value}
}

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
func (b *LigmaBoolean) MapKey() MapKey {
	var key MapKey
	if b.Value {
		key.Value = 1
	} else {
		key.Value = 0
	}
	return MapKey{Type: b.Type(), Value: key.Value}
}

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
func (s *LigmaString) MapKey() MapKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return MapKey{Type: s.Type(), Value: h.Sum64()}
}

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


type MapKey struct {
	Type ObjectType
	Value uint64
}


type MapPair struct {
	Key   LigmaObject
	Value LigmaObject
}

// LigmaMap
type LigmaMap struct {
	Pairs map[MapKey]MapPair
}

func (m *LigmaMap) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range m.Pairs {
		pairs = append(pairs, pair.Key.Inspect()+":"+pair.Value.Inspect())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

func (m *LigmaMap) Type() ObjectType { return MAP_OBJ }