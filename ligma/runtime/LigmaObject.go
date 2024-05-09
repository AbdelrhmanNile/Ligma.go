package runtime

type ObjectType string
type BuiltinLigmaFunction func(args ...LigmaObject) LigmaObject

const (
	INTEGER_OBJ = "INTEGER"
	FLOAT_OBJ = "FLOAT"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ = "NULL"
	STRING_OBJ = "STRING"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ = "ERROR"
	FUNCTION_OBJ = "FUNCTION"
	BUILTIN_OBJ = "BUILTIN"
	LIST_OBJ = "LIST"
	CLASS_OBJ = "CLASS"
	INSTANCE_OBJ = "INSTANCE"
)

type LigmaObject interface {
	Type() ObjectType
	Inspect() string
}

type LigmaCallable interface {
	LigmaObject
	Call(interpreter *Interpreter, args ...LigmaObject) LigmaObject
	Arity() int
}


type Error struct {
	Message string
}

func (e *Error) Inspect() string { return "ERROR: " + e.Message }
func (e *Error) Type() ObjectType { return ERROR_OBJ }
