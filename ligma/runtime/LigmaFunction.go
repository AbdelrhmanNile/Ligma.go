package runtime

import (
	"bytes"
	"fmt"
	"strings"
)

// LigmaFunction
type LigmaFunction struct {
	LigmaCallable
	Parameters []*Identifier
	Body *BlockStatement
	Env *Environment
}

func (f *LigmaFunction) Call(i *Interpreter, args ...LigmaObject) LigmaObject {
	env := NewEnclosedEnvironment(f.Env)

	for i, param := range f.Parameters {
		env.Set(param.Value, args[i])
	}

	return unwrapReturnValue(i.ExecuteBlock(f.Body, env))
}

func (f *LigmaFunction) Bind(instance *LigmaInstance) *LigmaFunction {
	env := NewEnclosedEnvironment(f.Env)
	env.Set("self", instance)
	return &LigmaFunction{Parameters: f.Parameters, Body: f.Body, Env: env}
	//return nil
}

func (f *LigmaFunction) Arity() int { return len(f.Parameters) }

func (f *LigmaFunction) Type() ObjectType { return FUNCTION_OBJ }
func (f *LigmaFunction) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("func")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

// ReturnValue
type ReturnValue struct {
	Value LigmaObject // fix this shit, reveiw the old code for return Unwrap
}

func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

// built-in functions
type Builtin struct {
	LigmaCallable
	Literal string
	Fn BuiltinLigmaFunction
	NumArgs int
}

func (b *Builtin) Call(i *Interpreter,args ...LigmaObject) LigmaObject { 
	return b.Fn(args...)
}
func (b *Builtin) Arity() int { return b.NumArgs }
func (b *Builtin) Inspect() string { return fmt.Sprintf("<built-in function %s>", b.Literal) }
func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
// bind 
func (b *Builtin) Bind(instance *LigmaInstance) *Builtin {
	return b
}


type BuiltinClassMethod struct {
	LigmaCallable
	Literal string
	Fn BuiltinLigmaFunction
	NumArgs int
	ObjInstance *LigmaInstance
}

func (b *BuiltinClassMethod ) Call(i *Interpreter, args ...LigmaObject) LigmaObject { 
	// extend args
	if b.Literal != "init" {
		args = append(args, b.ObjInstance)
	}
	return b.Fn(args...)
}
func (b *BuiltinClassMethod)  Arity() int { return b.NumArgs }
func (b *BuiltinClassMethod) Inspect() string { return fmt.Sprintf("<built-in function %s>", b.Literal) }
func (b *BuiltinClassMethod) Type() ObjectType { return BUILTIN_OBJ }
// bind 
func (b *BuiltinClassMethod) Bind(instance *LigmaInstance) *BuiltinClassMethod {
	b.ObjInstance = instance
	return b
}

type ClassMethods struct {
	BuiltinMethods map[string]*BuiltinClassMethod
	UserDefinedMethods map[string]*LigmaFunction
}
