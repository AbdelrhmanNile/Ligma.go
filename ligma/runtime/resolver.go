package runtime

import "os"

// Function types
const (
	_ int = iota
	ft_NONE
	ft_FUNCTION
	ft_INITIALIZER
	ft_METHOD
)

// Class types
const (
	_ int = iota
	cls_NONE
	cls_CLASS
	cls_SUBCLASS
)

const (
	_ int = iota
	loop_NONE
	loop_WHILE
)


type Resolver struct {
	interpreter *Interpreter
	// scopes is a stack of maps, where each map represents a scope
	scopes []map[string]bool

	currentFunction int
	currentClass int
	currentLoop int
}

func NewResolver(interpreter *Interpreter) *Resolver {
	currentFunction := ft_NONE
	currentClass := cls_NONE
	currentLoop := loop_NONE
	return &Resolver{interpreter: interpreter, scopes: []map[string]bool{}, currentFunction: currentFunction, currentClass: currentClass, currentLoop: currentLoop}
}

func (r *Resolver) Resolve(stmts []Statement) {
	for _, stmt := range stmts {
		r.resolveStatement(stmt)
	}
}

func (r *Resolver) resolveStatement(stmt Statement) {
	stmt.Accept(r)
}

func (r *Resolver) resolveExpression(expr Expression) {
	expr.Accept(r)
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) declare(name *Identifier) {
	if len(r.scopes) == 0 {
		return
	}

	if _, ok := r.scopes[len(r.scopes)-1][name.Value]; ok {
		NewError("Variable with this name already declared in this scope.")
		println("Variable with this name already declared in this scope.")
		println("Please come up with a better error handling mechanism.")
		os.Exit(1)
	}

	// print scope
	/* for k, v := range r.scopes[len(r.scopes)-1] {
	 	println(k, v)
	 } */

	r.scopes[len(r.scopes)-1][name.Value] = false
}

func (r *Resolver) define(name *Identifier) {
	if len(r.scopes) == 0 {
		return
	}

	r.scopes[len(r.scopes)-1][name.Value] = true
}

func (r *Resolver) resolveLocal(expr Expression, name string) {

	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name]; ok {
			r.interpreter.Resolve(expr, len(r.scopes)-1-i)
			return 
		}
	}

}

func (r *Resolver) VisitBlockStatement(block *BlockStatement) LigmaObject {
	r.beginScope()
	r.Resolve(block.Statements)
	r.endScope()
	return nil
}

func (r *Resolver) VisitDefStatement(def *DefStatement) LigmaObject {
	r.declare(def.Name)

	if def.Value != nil {
		r.resolveExpression(def.Value)
	}

	r.define(def.Name)
	return nil
}

func (r *Resolver) VisitExpressionStatement(exprStmt *ExpressionStatement) LigmaObject {
	r.resolveExpression(exprStmt.Expression)
	return nil
}

func (r *Resolver) VisitAssignExpression(assign *AssignExpression) LigmaObject {
	r.resolveExpression(assign.Value)
	r.resolveLocal(assign, assign.Name.Value)
	return nil
}

func (r *Resolver) VisitFunctionLiteral(funcLit *FunctionLiteral) LigmaObject {

	r.resolveFunction(funcLit, ft_FUNCTION)
	return nil
}

func (r *Resolver) resolveFunction(funcLit *FunctionLiteral, functionType int) {
	enclosingFunction := r.currentFunction
	r.currentFunction = functionType

	r.beginScope()
	for _, param := range funcLit.Parameters {
		r.declare(param)
		r.define(param)
	}

	r.Resolve(funcLit.Body.Statements)
	r.endScope()

	r.currentFunction = enclosingFunction
}


func (r *Resolver) VisitIdentifier(ident *Identifier) LigmaObject {

	if (len(r.scopes) > 0 && !r.scopes[len(r.scopes)-1][ident.Value]) {
		// if its a built-in function, don't resolve it
		if _, ok := builtins[ident.Value]; ok {
			return nil
		}

		if _, ok := builtinsClasses[ident.Value]; ok {
			return nil
		}

		// if we are in a class method, don't resolve it
		if r.currentFunction == ft_METHOD {
			r.resolveLocal(ident, ident.Value)
			return nil
		}

		if r.currentLoop == loop_WHILE {
			r.resolveLocal(ident, ident.Value)
			return nil
		}

		NewError("Can't read local variable in its own initializer.")
		println("Can't read local variable in its own initializer.")
		println("Please come up with a better error handling mechanism.")
		os.Exit(1)
	}

	r.resolveLocal(ident, ident.Value)

	return nil
}

func (r *Resolver) VisitIfExpression(ifExpr *IfExpression) LigmaObject {
	r.resolveExpression(ifExpr.Condition)
	r.resolveStatement(ifExpr.Consequence)

	if ifExpr.Alternative != nil {
		r.resolveStatement(ifExpr.Alternative)
	}

	return nil
}

func (r *Resolver) VisitReturnStatement(rs *ReturnStatement) LigmaObject {
	
	if r.currentFunction == ft_NONE {
		NewError("Can't return from top-level code.")
		println("Can't return from top-level code.")
		println("Please come up with a better error handling mechanism.")
		os.Exit(1)
	}

	if r.currentFunction == ft_INITIALIZER {
		NewError("Can't return a value from an initializer.")
		println("Can't return a value from an initializer.")
		println("Please come up with a better error handling mechanism.")
		os.Exit(1)
	}
	
	if rs.ReturnValue != nil {
		r.resolveExpression(rs.ReturnValue)
	}

	return nil
}

func (r *Resolver) VisitWhileStatement(ws *WhileStatement) LigmaObject {
	r.currentLoop = loop_WHILE
	r.resolveExpression(ws.Condition)
	r.resolveStatement(ws.Body)
	return nil
}

func (r *Resolver) VisitInfixExpression(ie *InfixExpression) LigmaObject {
	r.resolveExpression(ie.Left)
	r.resolveExpression(ie.Right)
	return nil
}

func (r *Resolver) VisitPrefixExpression(pe *PrefixExpression) LigmaObject {
	r.resolveExpression(pe.Right)
	return nil
}

func (r *Resolver) VisitCallExpression(ce *CallExpression) LigmaObject {
	r.resolveExpression(ce.Function)
	for _, arg := range ce.Arguments {
		r.resolveExpression(arg)
	}
	return nil
}

func (r *Resolver) VisitIntegerLiteral(il *IntegerLiteral) LigmaObject {
	return nil
}

func (r *Resolver) VisitFloatLiteral(fl *FloatLiteral) LigmaObject {
	return nil
}

func (r *Resolver) VisitBoolean(b *Boolean) LigmaObject {
	return nil
}

func (r *Resolver) VisitNull(n *Null) LigmaObject {
	return nil
}

func (r *Resolver) VisitListLiteral(ll *ListLiteral) LigmaObject {
	for _, elem := range ll.Elements {
		r.resolveExpression(elem)
	}
	return nil
}

func (r *Resolver) VisitStringLiteral(sl *StringLiteral) LigmaObject {
	return nil
}

func (r *Resolver) VisitIndexExpression(ie *IndexExpression) LigmaObject {
	r.resolveExpression(ie.Left)
	r.resolveExpression(ie.Index)
	return nil
}

func (r *Resolver) VisitGetExpression(ge *GetExpression) LigmaObject {
	r.resolveExpression(ge.Object)
	return nil
}

func (r *Resolver) VisitSetExpression(se *SetExpression) LigmaObject {
	r.resolveExpression(se.Value)
	r.resolveExpression(se.Object)
	return nil
}

func (r *Resolver) VisitSelfExpression(se *Self) LigmaObject {

	if r.currentClass == cls_NONE {
		NewError("Can't use 'self' outside of a class.")
		println("Can't use 'self' outside of a class.")
		println("Please come up with a better error handling mechanism.")
		os.Exit(1)
	}

	r.resolveLocal(se, se.Token.Literal)
	return nil
}

func (r *Resolver) VisitClassStatement(cs *Class) LigmaObject {

	enclosingClass := r.currentClass
	r.currentClass	= cls_CLASS

	r.declare(cs.Name)
	r.define(cs.Name)

	if (cs.Superclass != nil ){
		if cs.Name.Value == cs.Superclass.Value {
			NewError("A class can't inherit from itself.")
			println("A class can't inherit from itself.")
			println("Please come up with a better error handling mechanism.")
			os.Exit(1)
		}
		r.currentClass = cls_SUBCLASS
		r.resolveExpression(cs.Superclass)
	}

	if cs.Superclass != nil {
		r.beginScope()
		r.scopes[len(r.scopes)-1]["super"] = true
	}

	r.beginScope()
	r.scopes[len(r.scopes)-1]["self"] = true
	for _, method := range cs.Methods {
		method_func := method.Value.(*FunctionLiteral)
		declarition := ft_METHOD
		if method.Name.Value == "init" {
			declarition = ft_INITIALIZER
		}
		r.resolveFunction(method_func, declarition)
	}
	r.endScope()


	if cs.Superclass != nil {
		r.endScope()
	}

	r.currentClass = enclosingClass

	return nil
}

func (r *Resolver) VisitSuper(se *Super) LigmaObject {

	if r.currentClass == cls_NONE {
		NewError("Can't use 'super' outside of a class.")
		println("Can't use 'super' outside of a class.")
		println("Please come up with a better error handling mechanism.")
		os.Exit(1)
	} else if r.currentClass != cls_SUBCLASS {
		NewError("Can't use 'super' in a class with no superclass.")
		println("Can't use 'super' in a class with no superclass.")
		println("Please come up with a better error handling mechanism.")
		os.Exit(1)
	}

	r.resolveLocal(se, se.Token.Literal)
	return nil
}