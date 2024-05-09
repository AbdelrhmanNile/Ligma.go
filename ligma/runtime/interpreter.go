package runtime

import (
	"fmt"
)

var (
	// NULL is the null object
	NULL = &LigmaNull{}
	// TRUE is the true object
	TRUE = &LigmaBoolean{Value: true}
	// FALSE is the false object
	FALSE = &LigmaBoolean{Value: false}
)

type TreeWalker interface {
	ExpressionVisitor
	StatementVisitor
	//VisitProgram(*Program) LigmaObject
}

type Interpreter struct {
	globals *Environment 
	locals map[Expression]int
	Env *Environment
}

func NewInterpreter() *Interpreter {
	globals := NewEnvironment()

	DefineBuiltinTypes()
	
	// add built-in functions
	for name, builtin := range builtins {
		globals.Set(name, builtin)
	}

	// add built-in classes
	for name, class := range builtinsClasses{
		globals.Set(name, class)
	}

	env := globals

	return &Interpreter{globals: globals, locals: make(map[Expression]int), Env: env}
}

func (i *Interpreter) Resolve(expr Expression, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) ExecuteStatement(statement Statement) LigmaObject {
	return statement.Accept(i)
}

func (i *Interpreter) EvaluateExpression(expression Expression) LigmaObject {
	return expression.Accept(i)
}


func (i *Interpreter) Interpret(p *Program) LigmaObject {
	var result LigmaObject

	for _, statement := range p.Statements {
		result = i.ExecuteStatement(statement)
	}

	return result
}

func (i *Interpreter) ExecuteBlock(block *BlockStatement, env *Environment) LigmaObject {

	previousEnv := i.Env
	i.Env = env



	for _, statement := range block.Statements {
		result := i.ExecuteStatement(statement)
		if result != nil {
			if result.Type() == RETURN_VALUE_OBJ || result.Type() == ERROR_OBJ {
				i.Env = previousEnv
				return result
			}
		}
	}


	i.Env = previousEnv

	return nil
}

func (i *Interpreter) LookupVariable(name string, expr Expression) LigmaObject {
	distance, ok := i.locals[expr]
	if ok {
		ret, _ := i.Env.GetAt(distance, name)
		return ret
	}
	ret , ok := i.globals.Get(name)
	
	if ok {
		return ret
	}

	return NewError("undefined variable %s", name)
}

func (i *Interpreter) VisitDefStatement(def *DefStatement) LigmaObject{
	val := def.Value.Accept(i)
	if isError(val) {
		return val
	}

	if _, ok := builtins[def.Name.Value]; ok {
		return NewError("Built-in function %s cannot be redefined", def.Name.Value)
	}


	i.Env.Set(def.Name.Value, val)
	return nil
}

func (i *Interpreter) VisitReturnStatement(rs *ReturnStatement) LigmaObject {
	val := rs.ReturnValue.Accept(i)
	if isError(val) {
		return val
	}
	return &ReturnValue{Value: val}
}

func (i *Interpreter) VisitBlockStatement(block *BlockStatement) LigmaObject {
	//var result LigmaObject

	/* for _, statement := range block.Statements {
		result = statement.Accept(i)

		if result != nil && (result.Type() == RETURN_VALUE_OBJ || result.Type() == ERROR_OBJ) {
			return result
		}
	}

 */	

	return i.ExecuteBlock(block, NewEnclosedEnvironment(i.Env))
	
}

func (i *Interpreter) VisitClassStatement(class *Class) LigmaObject {

	i.Env.Set(class.Name.Value, nil)


	classObj := &LigmaClass{Name: class.Name.Value}

	objectClass := builtinsClasses["object"]

	classObj.Superclasses = append(classObj.Superclasses, objectClass)


	if class.Superclass != nil {
		superclass := class.Superclass.Accept(i)
		if isError(superclass) {
			return superclass
		}

		if superclass.Type() != CLASS_OBJ {
			return NewError("superclass must be a class")
		}

		//classObj.Superclass = superclass.(*LigmaClass)	
		classObj.Superclasses = append(classObj.Superclasses, superclass.(*LigmaClass))

		
	}

	/* if class.Superclass != nil {
		i.Env = NewEnclosedEnvironment(i.Env)
		i.Env.Set("super", classObj.Superclass)
	} */

	if class.Superclass == nil {
		// set super to Object
		i.Env = NewEnclosedEnvironment(i.Env)
		i.Env.Set("super", objectClass)
	} else {
		// inherit from Object then from the superclass
		i.Env = NewEnclosedEnvironment(i.Env)
		i.Env.Set("super", objectClass)
		i.Env = NewEnclosedEnvironment(i.Env)
		i.Env.Set("super", classObj.Superclasses[1])
	}

	
	methods := make(map[string]*LigmaFunction)


	for _, method := range class.Methods {
		method_func := method.Value.(*FunctionLiteral)
		methods[method.Name.Value] = &LigmaFunction{Parameters: method_func.Parameters, Body: method_func.Body, Env: i.Env}
	}

	classObj.Methods = ClassMethods{UserDefinedMethods: methods}

	/* if class.Superclass != nil {
		i.Env = i.Env.parent
	} */

	if class.Superclass != nil {
		i.Env = i.Env.parent
		i.Env = i.Env.parent
	} else {
		i.Env = i.Env.parent
	}
	
	i.Env.Set(class.Name.Value, classObj)

	return nil
}

func (i *Interpreter) VisitSuper(s *Super) LigmaObject {
	distance, ok := i.locals[s]
	if !ok {
		return NewError("super must be used inside a method")
	}

	superclass, ok := i.Env.GetAt(distance, "super")
	if !ok {
		return NewError("super must be used inside a method")
	}

	supercls := superclass.(*LigmaClass)

	thisInstance, ok := i.Env.GetAt(distance - 1, "self")
	if !ok {
		return NewError("super must be used inside a method")
	}

	thisInstanceObj := thisInstance.(*LigmaInstance)

	method := supercls.GetMethod(s.Method.Value)

	if method == nil {
		return NewError("undefined method %s", s.Method.Value)
	}

	//return method.Bind(thisInstanceObj)

	if method.UserMethod != nil {
		return method.UserMethod.Bind(thisInstanceObj)
	}

	return method.BuiltinMethod.Bind(thisInstanceObj)
}

func (i *Interpreter) VisitWhileStatement(ws *WhileStatement) LigmaObject {
	condition := i.EvaluateExpression(ws.Condition)
	if isError(condition) {
		return condition
	}

	for isTruthy(condition) {
		result := i.ExecuteStatement(ws.Body)
		if isError(result) {
			return result
		}
		condition = i.EvaluateExpression(ws.Condition)
	}
	return nil
}

func (i *Interpreter) VisitExpressionStatement(es *ExpressionStatement) LigmaObject {
	return es.Expression.Accept(i)
}

// literals 
func (i *Interpreter) VisitIntegerLiteral(il *IntegerLiteral) LigmaObject {
	//return &LigmaInteger{Value: il.Value}
	int_class, _ := i.Env.Get("int")
	return ApplyFunction(i, int_class.(*LigmaClass), []LigmaObject{&LigmaInteger{Value: il.Value}})
}

func (i *Interpreter) VisitFloatLiteral(fl *FloatLiteral) LigmaObject {
	//return &LigmaFloat{Value: fl.Value}
	float_class, _ := i.Env.Get("float")
	return ApplyFunction(i, float_class.(*LigmaClass), []LigmaObject{&LigmaFloat{Value: fl.Value}})
}

func (i *Interpreter) VisitBoolean(b *Boolean) LigmaObject {
	return &LigmaBoolean{Value: b.Value}
}

func (i *Interpreter) VisitNull(n *Null) LigmaObject {
	return NULL
}

func (i *Interpreter) VisitListLiteral(ll *ListLiteral) LigmaObject {
	elements := LigmaList{}

	for _, element := range ll.Elements {
		//elements = append(elements, element.Accept(i))
		elements.Elements = append(elements.Elements, element.Accept(i))
	}

	//return &LigmaList{Elements: elements}
	list_class, _ := i.Env.Get("list")
	return ApplyFunction(i, list_class.(*LigmaClass), []LigmaObject{&elements})
}

func (i *Interpreter) VisitStringLiteral(sl *StringLiteral) LigmaObject {
	//return &LigmaString{Value: sl.Value}
	string_class, _ := i.Env.Get("str")
	return ApplyFunction(i, string_class.(*LigmaClass), []LigmaObject{&LigmaString{Value: sl.Value}})

}

func (i *Interpreter) VisitFunctionLiteral(fl *FunctionLiteral) LigmaObject {
	return &LigmaFunction{Parameters: fl.Parameters, Body: fl.Body, Env: i.Env}
}

func (i *Interpreter) VisitPrefixExpression(pe *PrefixExpression) LigmaObject {
	right := pe.Right.Accept(i)
	if isError(right) {
		return right
	}
	operator := pe.Operator

	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
}
	return NewError("unknown operator: %s%s", operator, right.Type())
}

func (i *Interpreter) VisitInfixExpression(ie *InfixExpression) LigmaObject {
	left := ie.Left.Accept(i)
	if isError(left) {
		return left
	}

	right := ie.Right.Accept(i)
	if isError(right) {
		return right
	}

	operator := ie.Operator

	switch {
		/* case left.Type() == ObjectType("int") && right.Type() == ObjectType("int"):
			return evalIntegerInfixExpression(i, operator, left, right)
		case left.Type() == FLOAT_OBJ && right.Type() == FLOAT_OBJ:
			return evalFloatInfixExpression(operator, left, right)
		case left.Type() == INTEGER_OBJ && right.Type() == FLOAT_OBJ:
			return evalMixedInfixExpression(operator, left, right)
		case left.Type() == FLOAT_OBJ && right.Type() == INTEGER_OBJ:
				return evalMixedInfixExpression(operator, left, right) */
		
		case operator == "+":
			return left.(*LigmaInstance).Add(right.(*LigmaInstance))
		case operator == "-":
			return left.(*LigmaInstance).Sub(right.(*LigmaInstance))
		case operator == "*":
			return left.(*LigmaInstance).Mul(right.(*LigmaInstance))
		case operator == "/":
			return left.(*LigmaInstance).Div(right.(*LigmaInstance))
		case operator == "%":
			return left.(*LigmaInstance).Mod(right.(*LigmaInstance))
		case operator == "<":
			return left.(*LigmaInstance).Lt(right.(*LigmaInstance))
		
		case operator == "==":
			return left.(*LigmaInstance).Eq(right.(*LigmaInstance))
		case operator == "!=":
			return left.(*LigmaInstance).Ne(right.(*LigmaInstance))
		case operator == "and":
			return nativeBoolToBooleanObject(isTruthy(left) && isTruthy(right))
		case operator == "or":
			return nativeBoolToBooleanObject(isTruthy(left) || isTruthy(right))
		
	}

	return NewError("unknown operator: %s %s %s", left.Type(), ie.Operator, right.Type())
}

func (i *Interpreter) VisitIfExpression(ie *IfExpression) LigmaObject {
	condition := i.EvaluateExpression(ie.Condition)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return i.ExecuteStatement(ie.Consequence)
	} else if ie.Alternative != nil {
		return i.ExecuteStatement(ie.Alternative)
	} else {
		return NULL
	}
}

func (i *Interpreter) VisitIdentifier(id *Identifier) LigmaObject {
	return i.LookupVariable(id.Value, id)
}

func (i *Interpreter) VisitIndexExpression(ie *IndexExpression) LigmaObject {
	left := ie.Left.Accept(i)
	if isError(left) {
		return left
	}

	index := ie.Index.Accept(i)
	if isError(index) {
		return index
	}
	
	get_func, ok := left.(*LigmaInstance).Get("__get__")

	if !ok {
		return NewError("object of type %s does not support indexing", left.Type())
	}

	return ApplyFunction(i, get_func, []LigmaObject{index})
}

func (i *Interpreter) VisitAssignExpression(ae *AssignExpression) LigmaObject {
	val := ae.Value.Accept(i)
	if isError(val) {
		return val
	}

	if _, ok := builtins[ae.Name.Value]; ok {
		return NewError("identifier %s is reserved", ae.Name.Value)
	}

	if distance, ok := i.locals[ae]; ok { // if the variable is local
		i.Env.SetAt(distance, ae.Name.Value, val)
	} else {
		i.globals.Set(ae.Name.Value, val)
	}

	return nil
}

func (i *Interpreter) VisitCallExpression(ce *CallExpression) LigmaObject {

	//os.Exit(1)
	function := ce.Function.Accept(i)
	if isError(function) {
		return function
	}
	
	args := []LigmaObject{}
	for _, arg := range ce.Arguments {
		evalArg := arg.Accept(i)
		if isError(evalArg) {
			return evalArg
		}
		args = append(args, evalArg)
	}
	
	return ApplyFunction(i, function, args)
}

func (i *Interpreter) VisitSelfExpression(se *Self) LigmaObject {
	/* instance, ok := i.Env.Get("self")
	if !ok {
		return NewError("no self object found")
	}
	return instance */

	return i.LookupVariable(se.Token.Literal, se)
}

func (i *Interpreter) VisitGetExpression(ge *GetExpression) LigmaObject {
	obj := ge.Object.Accept(i)
	if isError(obj) {
		return obj
	}
	return evalGetExpression(obj, ge.Property)
}

func (i *Interpreter) VisitSetExpression(se *SetExpression) LigmaObject {
	obj := se.Object.Accept(i)
	if isError(obj) {
		return obj
	}

	val := se.Value.Accept(i)
	if isError(val) {
		return val
	}

	//obj.(*LigmaInstance).Set(se.Property.Value, val)
	//obj.(*BaseObjectInstance).Set(se.Property.Value, val)

	// is it LigmaInstance or BaseObjectInstance?
	switch obj := obj.(type) {
		case *LigmaInstance:
			obj.Set(se.Property.Value, val)
		//case *BaseObjectInstance:
		//	obj.Set(se.Property.Value, val)
	}
	return val
}

func (i *Interpreter) createIntInstance(number *LigmaInteger) *LigmaInstance {
	intClass, _ := i.Env.Get("int")
	instance := &LigmaInstance{Class: intClass.(*LigmaClass), Fields: map[string]LigmaObject{
		"value": number,
	}}
	return instance

}

func isError(obj LigmaObject) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}

func NewError(format string, a ...interface{}) *Error {
	println(fmt.Sprintf(format, a...))
	//os.Exit(1)
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func nativeBoolToBooleanObject(input bool) *LigmaBoolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalBangOperatorExpression(right LigmaObject) LigmaObject {
	switch right {
		case TRUE:
			return FALSE
		case FALSE:
			return TRUE
		case NULL:
			return TRUE
	}
	return FALSE
}

func evalMinusPrefixOperatorExpression(right LigmaObject) LigmaObject {
	switch right := right.(type) {
		case *LigmaInteger:
			return &LigmaInteger{Value: -right.Value}
		case *LigmaFloat:
			return &LigmaFloat{Value: -right.Value}
		
	}
	return NewError("unknown operator: -%s", right.Type())
}


func evalIndexExpression(left, index LigmaObject) LigmaObject {
	switch {
	case left.Type() == LIST_OBJ && index.Type() == ObjectType("int"):
			return evalListIndexExpression(left, index)
	}

	return NewError("index operator not supported: %s", left.Type())
}

func evalListIndexExpression(list, index LigmaObject) LigmaObject {
	listObject := list.(*LigmaInstance).Fields["value"].(*LigmaList)
	
	//idx := index.(*LigmaInteger).Value
	idx := index.(*LigmaInstance).Fields["value"].(*LigmaInteger).Value

	max := int64(len(listObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return listObject.Elements[idx]
}

func ApplyFunction(i *Interpreter, fn LigmaObject, args []LigmaObject) LigmaObject {
	function, ok := fn.(LigmaCallable)

		if !ok {
			return NewError("not a function: %s", fn.Type())
		} else {
			if function.Arity() != -1 {
			if len(args) != function.Arity() {
				return NewError("wrong number of arguments. got=%d, want=%d", len(args), function.Arity())
			}
		}
		
	}
	return function.Call(i, args...)
}

func evalGetExpression(obj LigmaObject, property *Identifier) LigmaObject {
	switch obj := obj.(type) {
		case *LigmaClass:
			/* if method, ok := obj.Methods[property.Value]; ok {
				return method */
			if method := obj.GetMethod(property.Value); method != nil {
				if method.UserMethod != nil {
					return method.UserMethod
				}
				return method.BuiltinMethod
			}

			return NewError("no method %s found for class %s", property.Value, obj.Name)
		case *LigmaInstance:
			//val, _ := obj.Get(property.Value)
			//return val
			val, _ := obj.Get(property.Value)
			
			switch val := val.(type) {
				case *LigmaFunction:
					return val
				case *BuiltinClassMethod:
					val.ObjInstance = obj
					return val
			}
			return val
			
		case *ReturnValue:
			return evalGetExpression(obj.Value, property)

/* 			switch val := val.(type) {
				case *LigmaFunction:
					return val
				case *BuiltinClassMethod:
					val.ObjInstance = obj
					return val

				} */
			//case *BaseObjectInstance:
		//	val, _ := obj.Get(property.Value)
		//	vall, _ := val.(*BuiltinClassMethod)
		//	vall.ObjInstance = obj
		//	return vall
	}
	return NewError("property access not supported on %s", obj.Type())
}


func isTruthy(obj LigmaObject) bool {
	switch obj {
		case NULL:
			return false
		case TRUE:
			return true
		case FALSE:
			return false
	}

	return true
}

func unwrapReturnValue(obj LigmaObject) LigmaObject {
	if returnValue, ok := obj.(*ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}