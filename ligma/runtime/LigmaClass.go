package runtime

import (
	"fmt"
)

type MethodWrapper struct {
	UserMethod *LigmaFunction
	BuiltinMethod *BuiltinClassMethod
}

type LigmaClass struct {
	Name string
	Superclasses []*LigmaClass
	//Methods map[string]*LigmaFunction
	Methods ClassMethods
}

func (c *LigmaClass) Call(i *Interpreter, args ...LigmaObject) LigmaObject {
	instance := &LigmaInstance{Class: c, Fields: map[string]LigmaObject{}, interpreter: i}

	/* constructor, ok := c.Methods["init"]
	if ok {
		(*constructor.Bind(instance)).Call(i, args...)
	}
	return instance */

	// check if the class has an init user defined method
	constructor, ok := c.Methods.UserDefinedMethods["init"]
	if ok {
		(*constructor).Bind(instance).Call(i, args...)
		return instance
	}

	// check if the class has an init builtin method
	constructorBuiltin, ok := c.Methods.BuiltinMethods["init"]
	if ok {
		args = append(args, instance)
		(*constructorBuiltin).Bind(instance).Call(i, args...)
		return instance
	}

	return instance
}

func (c *LigmaClass) GetMethod(name string) *MethodWrapper {
	/* method, ok := c.Methods[name]
	if ok {
		return method
	} */


	// check if the method is user defined
	method, ok := c.Methods.UserDefinedMethods[name]
	if ok {
		return &MethodWrapper{UserMethod: method}
	}

	// check if the method is builtin
	methodBuiltin, ok := c.Methods.BuiltinMethods[name]
	if ok {
		return &MethodWrapper{BuiltinMethod: methodBuiltin}
	}

	/* if c.Superclass != nil {
		return c.Superclass.GetMethod(name)
	} */

	if len(c.Superclasses) > 0 {
		// loop backwards through the superclasses
		for i := len(c.Superclasses) - 1; i >= 0; i-- {
			method := c.Superclasses[i].GetMethod(name)
			if method != nil {
				return method
			}
		}
	}
	return nil
}

func (c *LigmaClass) Arity() int { 

	// check if the class has an init user defined method
	constructor, ok := c.Methods.UserDefinedMethods["init"]
	if ok {
		return constructor.Arity()
	}
	
	// check if the class has an init builtin method
	constructorBuiltin, ok := c.Methods.BuiltinMethods["init"]
	if ok {
		return constructorBuiltin.Arity()
	}

	return 0
}

func (c *LigmaClass) Inspect() string {
	return fmt.Sprintf("<class %s>", c.Name)
}

func (c *LigmaClass) Type() ObjectType { return CLASS_OBJ }

type LigmaInstance struct {
	Class *LigmaClass
	Fields map[string]LigmaObject

	interpreter *Interpreter
}

func (i *LigmaInstance) Inspect() string {
	return fmt.Sprintf("<instance of %s>", i.Class.Name)
}

func (i *LigmaInstance) Type() ObjectType { return ObjectType(i.Class.Name) }

func (i *LigmaInstance) Get(name string) (LigmaObject, bool) {
	field, ok := i.Fields[name]
	if ok {
		return field, true
	}

	method := i.Class.GetMethod(name)
	if method != nil {
		//return method.Bind(i), true
		if method.UserMethod != nil {
			return method.UserMethod.Bind(i), true
		}
		if method.BuiltinMethod != nil {
			return method.BuiltinMethod.Bind(i), true
		}
	}

	return &LigmaNull{}, false
}

func (i *LigmaInstance) Set(name string, val LigmaObject) LigmaObject {
	i.Fields[name] = val
	return val
}

func (i *LigmaInstance) Add (other LigmaObject) LigmaObject {
	add_func, _ := i.Get("__add__")
	switch add_func.(type) {
	case *LigmaFunction:
		return add_func.(*LigmaFunction).Bind(i).Call(nil, other)
	case *BuiltinClassMethod:
		return add_func.(*BuiltinClassMethod).Bind(i).Call(nil, other)
	}
	return &LigmaNull{}
}

func (i *LigmaInstance) Sub (other LigmaObject) LigmaObject {
	sub_func, _ := i.Get("__sub__")
	switch sub_func.(type) {
	case *LigmaFunction:
		return sub_func.(*LigmaFunction).Bind(i).Call(nil, other)
	case *BuiltinClassMethod:
		return sub_func.(*BuiltinClassMethod).Bind(i).Call(nil, other)
	}
	return &LigmaNull{}
}

func (i *LigmaInstance) Eq (other LigmaObject) LigmaObject {
	eq_func, _ := i.Get("__eq__")
	switch eq_func.(type) {
	case *LigmaFunction:
		return eq_func.(*LigmaFunction).Bind(i).Call(nil, other)
	case *BuiltinClassMethod:
		return eq_func.(*BuiltinClassMethod).Bind(i).Call(nil, other)
	}
	return &LigmaNull{}
}

func (i *LigmaInstance) Ne (other LigmaObject) LigmaObject {
	ne_func, _ := i.Get("__ne__")
	switch ne_func.(type) {
	case *LigmaFunction:
		return ne_func.(*LigmaFunction).Bind(i).Call(nil, other)
	case *BuiltinClassMethod:
		return ne_func.(*BuiltinClassMethod).Bind(i).Call(nil, other)
	}
	return &LigmaNull{}
}

func (i *LigmaInstance) Mul (other LigmaObject) LigmaObject {
	mul_func, _ := i.Get("__mul__")
	switch mul_func.(type) {
	case *LigmaFunction:
		return mul_func.(*LigmaFunction).Bind(i).Call(nil, other)
	case *BuiltinClassMethod:
		return mul_func.(*BuiltinClassMethod).Bind(i).Call(nil, other)
	}
	return &LigmaNull{}
}

func (i *LigmaInstance) Div (other LigmaObject) LigmaObject {
	div_func, _ := i.Get("__div__")
	switch div_func.(type) {
	case *LigmaFunction:
		return div_func.(*LigmaFunction).Bind(i).Call(nil, other)
	case *BuiltinClassMethod:
		return div_func.(*BuiltinClassMethod).Bind(i).Call(nil, other)
	}
	return &LigmaNull{}
}

func (i *LigmaInstance) Mod (other LigmaObject) LigmaObject {
	mod_func, _ := i.Get("__mod__")
	switch mod_func.(type) {
	case *LigmaFunction:
		return mod_func.(*LigmaFunction).Bind(i).Call(nil, other).(LigmaObject)
	case *BuiltinClassMethod:
		return mod_func.(*BuiltinClassMethod).Bind(i).Call(nil, other).(LigmaObject)
	}
	return &LigmaNull{}
}

func (i *LigmaInstance) Lt (other LigmaObject) LigmaObject {
	lt_func, _ := i.Get("__lt__")
	switch lt_func.(type) {
	case *LigmaFunction:
		return lt_func.(*LigmaFunction).Bind(i).Call(nil, other)
	case *BuiltinClassMethod:
		return lt_func.(*BuiltinClassMethod).Bind(i).Call(nil, other)
	}
	return &LigmaNull{}
}