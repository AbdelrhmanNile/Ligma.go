package runtime

import (
	"fmt"
	"strings"
	"time"
)

var builtins = map[string]*Builtin{
"len": {
		Literal: "len",
		Fn: func(args ...LigmaObject) LigmaObject {
			obj := args[0]
			__len__, ok := obj.(*LigmaInstance).Get("__len__")
			if ok {
				return __len__.(LigmaCallable).Call(nil, obj)
			}
			return NewError("object of type '%s' has no len()", obj.Type())
		},
		NumArgs: 1,
	},

	"print": {
		Literal: "print",
		Fn: func(args ...LigmaObject) LigmaObject {
			for _, arg := range args {
				// assert that the argument is a LigmaInstance
				if instance, ok := arg.(*LigmaInstance); ok {
					// check if the instance has a __str__ method
					str, ok := instance.Get("__str__")
					if ok {
						
						// call the __str__ method
						switch str := str.(type) {
						case *BuiltinClassMethod:
							println(str.Call(nil, instance).Inspect())
						case *LigmaFunction:
							res := str.Call(instance.interpreter, instance)
							res = res.(*LigmaInstance).Fields["value"]
							println(res.Inspect())
						}
					}
				}
			}
			return  nil
		},
		NumArgs: -1,
	},

	"time": {
		Literal: "time",
		Fn: func(args ...LigmaObject) LigmaObject {
			t := time.Now().UnixMilli()
			return builtinsClasses["int"].Call(nil, &LigmaInteger{Value: t})
		},
	},

	"input": {
		Literal: "input",
		Fn: func(args ...LigmaObject) LigmaObject {
			print(args[0].(*LigmaInstance).Fields["value"].(*LigmaString).Value	)
			var input string
			_, err := fmt.Scanln(&input)
			if err != nil {
				return NewError(err.Error())
			}
			return builtinsClasses["str"].Call(nil, &LigmaString{Value: input})
		},
		NumArgs: 1,
	},


}



 var builtinsClasses = map[string]*LigmaClass{
	"object": {
		Name: "object",
		Methods: ClassMethods{
			BuiltinMethods: map[string]*BuiltinClassMethod{
				"__repr__": {
					Literal: "__repr__",
					Fn: func(args ...LigmaObject) LigmaObject {
						//self := args[len(args)-1].(*LigmaInstance)
						return &LigmaString{Value: "<object instance>"}
					},
				},
				"__str__": {
					Literal: "__str__",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						repr, _ := self.Get("__repr__")
						return repr.(*BuiltinClassMethod).Call(nil, self)
					},
				},
				"__eq__": {
					Literal: "__eq__",
					Fn: func(args ...LigmaObject) LigmaObject {
						return NewError("Not implemented")
					},
					NumArgs: 1,
				},
				"__ne__": {
					Literal: "__ne__",
					Fn: func(args ...LigmaObject) LigmaObject {
						return NewError("Not implemented")
					},
					NumArgs: 1,
				},
				"__add__": {
					Literal: "__add__",
					Fn: func(args ...LigmaObject) LigmaObject {
						return NewError("Not implemented")
					},
					NumArgs: 1,
				},
				"__sub__": {
					Literal: "__sub__",
					Fn: func(args ...LigmaObject) LigmaObject {
						return NewError("Not implemented")
					},
					NumArgs: 1,
				},
				"__mul__": {
					Literal: "__mul__",
					Fn: func(args ...LigmaObject) LigmaObject {
						return NewError("Not implemented")
					},
					NumArgs: 1,
				},
				"__div__": {
					Literal: "__div__",
					Fn: func(args ...LigmaObject) LigmaObject {
						return NewError("Not implemented")
					},
					NumArgs: 1,
				},
				"__mod__": {
					Literal: "__mod__",
					Fn: func(args ...LigmaObject) LigmaObject {
						return NewError("Not implemented")
					},
					NumArgs: 1,
				},	
			},

			UserDefinedMethods: map[string]*LigmaFunction{},
		},
	},
}

func DefineBuiltinTypes() {
	// implement a metaclass for the built-in classes, metaclass instances are classes themselves
	builtinsClasses["type"] = &LigmaClass{
		Name: "type",
		Methods: ClassMethods{
			BuiltinMethods: map[string]*BuiltinClassMethod{
				"__repr__": {
					Literal: "__repr__",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						return &LigmaString{Value: self.Fields["value"].(*LigmaString).Value}
					},
				},
				"__str__": {
					Literal: "__str__",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						return &LigmaString{Value: self.Fields["value"].(*LigmaString).Value}
					},
				},
			},
		},
		Superclasses: []*LigmaClass{builtinsClasses["object"]},
	}


	builtinsClasses["Number"] = &LigmaClass{
		Name: "Number",
		Methods: ClassMethods{
			BuiltinMethods: map[string]*BuiltinClassMethod{
				"__add__": {
					Literal: "__add__",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						args = args[:len(args)-1]

						my_type := self.Class.Name

						other := args[0].(*LigmaInstance)
						other_type := other.Class.Name

						switch my_type {
							case "int":
								switch other_type {
									case "int":
										my_val := self.Fields["value"].(*LigmaInteger).Value
										other_val := other.Fields["value"].(*LigmaInteger).Value
										return builtinsClasses["int"].Call(nil, &LigmaInteger{Value: my_val + other_val})
									case "float":
										my_val := float64(self.Fields["value"].(*LigmaInteger).Value)
										other_val := other.Fields["value"].(*LigmaFloat).Value
										return builtinsClasses["float"].Call(nil, &LigmaFloat{Value: my_val + other_val})
									}
							case "float":
								switch other_type {
									case "int":
										my_val := self.Fields["value"].(*LigmaFloat).Value
										other_val := float64(other.Fields["value"].(*LigmaInteger).Value)
										return builtinsClasses["float"].Call(nil, &LigmaFloat{Value: my_val + other_val})
									case "float":
										my_val := self.Fields["value"].(*LigmaFloat).Value
										other_val := other.Fields["value"].(*LigmaFloat).Value
										return builtinsClasses["float"].Call(nil, &LigmaFloat{Value: my_val + other_val})
									}
						}
						return NewError("unsupported operand type(s) for +: '%s' and '%s'", my_type, other_type)
					},
				},
				"__sub__": {
					Literal: "__sub__",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						args = args[:len(args)-1]

						my_type := self.Class.Name

						other := args[0].(*LigmaInstance)
						other_type := other.Class.Name

						switch my_type {
							case "int":
								switch other_type {
									case "int":
										my_val := self.Fields["value"].(*LigmaInteger).Value
										other_val := other.Fields["value"].(*LigmaInteger).Value
										return builtinsClasses["int"].Call(nil, &LigmaInteger{Value: my_val - other_val})
									case "float":
										my_val := float64(self.Fields["value"].(*LigmaInteger).Value)
										other_val := other.Fields["value"].(*LigmaFloat).Value
										return builtinsClasses["float"].Call(nil, &LigmaFloat{Value: my_val - other_val})
									}
							case "float":
								switch other_type {
									case "int":
										my_val := self.Fields["value"].(*LigmaFloat).Value
										other_val := float64(other.Fields["value"].(*LigmaInteger).Value)
										return builtinsClasses["float"].Call(nil, &LigmaFloat{Value: my_val - other_val})
									case "float":
										my_val := self.Fields["value"].(*LigmaFloat).Value
										other_val := other.Fields["value"].(*LigmaFloat).Value
										return builtinsClasses["float"].Call(nil, &LigmaFloat{Value: my_val - other_val})
									}
						}
						return NewError("unsupported operand type(s) for -: '%s' and '%s'", my_type, other_type)
					},
				},
				"__mul__": {
					Literal: "__mul__",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						args = args[:len(args)-1]

						my_type := self.Class.Name

						other := args[0].(*LigmaInstance)
						other_type := other.Class.Name
						

						switch my_type {
							case "int":
								switch other_type {
									case "int":
										my_val := self.Fields["value"].(*LigmaInteger).Value
										other_val := other.Fields["value"].(*LigmaInteger).Value
										return builtinsClasses["int"].Call(nil, &LigmaInteger{Value: my_val * other_val})
									case "float":
										my_val := float64(self.Fields["value"].(*LigmaInteger).Value)
										other_val := other.Fields["value"].(*LigmaFloat).Value
										return builtinsClasses["float"].Call(nil, &LigmaFloat{Value: my_val * other_val})
								}
							case "float":
								switch other_type {
									case "int":
										my_val := self.Fields["value"].(*LigmaFloat).Value
										other_val := float64(other.Fields["value"].(*LigmaInteger).Value)
										return builtinsClasses["float"].Call(nil, &LigmaFloat{Value: my_val * other_val})
									case "float":
										my_val := self.Fields["value"].(*LigmaFloat).Value
										other_val := other.Fields["value"].(*LigmaFloat).Value
										return builtinsClasses["float"].Call(nil, &LigmaFloat{Value: my_val * other_val})
								}
						}
						return NewError("unsupported operand type(s) for *: '%s' and '%s'", my_type, other_type)
					},
				},
				"__div__": {
					Literal: "__div__",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						args = args[:len(args)-1]

						my_type := self.Class.Name

						other := args[0].(*LigmaInstance)
						other_type := other.Class.Name

						switch my_type {
							case "int":
								switch other_type {
									case "int":
										my_val := self.Fields["value"].(*LigmaInteger).Value
										other_val := other.Fields["value"].(*LigmaInteger).Value
										return builtinsClasses["float"].Call(nil, &LigmaFloat{Value: float64(my_val) / float64(other_val)})
									case "float":
										my_val := float64(self.Fields["value"].(*LigmaInteger).Value)
										other_val := other.Fields["value"].(*LigmaFloat).Value
										return builtinsClasses["float"].Call(nil, &LigmaFloat{Value: my_val / other_val})
								}
							case "float":
								switch other_type {
									case "int":
										my_val := self.Fields["value"].(*LigmaFloat).Value
										other_val := float64(other.Fields["value"].(*LigmaInteger).Value)
										return builtinsClasses["float"].Call(nil, &LigmaFloat{Value: my_val / other_val})
									case "float":
										my_val := self.Fields["value"].(*LigmaFloat).Value
										other_val := other.Fields["value"].(*LigmaFloat).Value
										return builtinsClasses["float"].Call(nil, &LigmaFloat{Value: my_val / other_val})
								}
						}
						return NewError("unsupported operand type(s) for /: '%s' and '%s'", my_type, other_type)
					},
				},
				"__mod__": {
					Literal: "__mod__",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						args = args[:len(args)-1]

						my_type := self.Class.Name

						other := args[0].(*LigmaInstance)
						other_type := other.Class.Name

						switch my_type {
							case "int":
								switch other_type {
									case "int":
										my_val := self.Fields["value"].(*LigmaInteger).Value
										other_val := other.Fields["value"].(*LigmaInteger).Value
										return builtinsClasses["int"].Call(nil, &LigmaInteger{Value: my_val % other_val})
									case "float":
										my_val := float64(self.Fields["value"].(*LigmaInteger).Value)
										other_val := other.Fields["value"].(*LigmaFloat).Value
										return builtinsClasses["float"].Call(nil, &LigmaFloat{Value: float64(int64(my_val) % int64(other_val))})
								}
							case "float":
								switch other_type {
									case "int":
										my_val := self.Fields["value"].(*LigmaFloat).Value
										other_val := float64(other.Fields["value"].(*LigmaInteger).Value)
										return builtinsClasses["float"].Call(nil, &LigmaFloat{Value: float64(int64(my_val) % int64(other_val))})
									case "float":
										my_val := self.Fields["value"].(*LigmaFloat).Value
										other_val := other.Fields["value"].(*LigmaFloat).Value
										return builtinsClasses["float"].Call(nil, &LigmaFloat{Value: float64(int64(my_val) % int64(other_val))})
								}
						}
						return NewError("unsupported operand type(s) for %: '%s' and '%s'", my_type, other_type)
					},
				},
				"__eq__": {
					Literal: "__eq__",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						args = args[:len(args)-1]

						my_type := self.Class.Name

						other := args[0].(*LigmaInstance)
						other_type := other.Class.Name

						switch my_type {
							case "int":
								switch other_type {
									case "int":
										my_val := self.Fields["value"].(*LigmaInteger).Value
										other_val := other.Fields["value"].(*LigmaInteger).Value
										return nativeBoolToBooleanObject(my_val == other_val)
									case "float":
										my_val := float64(self.Fields["value"].(*LigmaInteger).Value)
										other_val := other.Fields["value"].(*LigmaFloat).Value
										return nativeBoolToBooleanObject(my_val == other_val)
								}
							case "float":
								switch other_type {
									case "int":
										my_val := self.Fields["value"].(*LigmaFloat).Value
										other_val := float64(other.Fields["value"].(*LigmaInteger).Value)
										return nativeBoolToBooleanObject(my_val == other_val)
									case "float":
										my_val := self.Fields["value"].(*LigmaFloat).Value
										other_val := other.Fields["value"].(*LigmaFloat).Value
										return nativeBoolToBooleanObject(my_val == other_val)
							}
						}
						return NewError("unsupported operand type(s) for ==: '%s' and '%s'", my_type, other_type)
					},
				},
				"__ne__": {
					Literal: "__ne__",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						args = args[:len(args)-1]

						my_type := self.Class.Name

						other := args[0].(*LigmaInstance)
						other_type := other.Class.Name

						switch my_type {
							case "int":
								switch other_type {
									case "int":
										my_val := self.Fields["value"].(*LigmaInteger).Value
										other_val := other.Fields["value"].(*LigmaInteger).Value
										return nativeBoolToBooleanObject(my_val != other_val)
									case "float":
										my_val := float64(self.Fields["value"].(*LigmaInteger).Value)
										other_val := other.Fields["value"].(*LigmaFloat).Value
										return nativeBoolToBooleanObject(my_val != other_val)
								}
							case "float":
								switch other_type {
									case "int":
										my_val := self.Fields["value"].(*LigmaFloat).Value
										other_val := float64(other.Fields["value"].(*LigmaInteger).Value)
										return nativeBoolToBooleanObject(my_val != other_val)
									case "float":
										my_val := self.Fields["value"].(*LigmaFloat).Value
										other_val := other.Fields["value"].(*LigmaFloat).Value
										return nativeBoolToBooleanObject(my_val != other_val)
							}
						}
						return NewError("unsupported operand type(s) for !=: '%s' and '%s'", my_type, other_type)
					},
				},
				"__lt__": {
					Literal: "__lt__",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						args = args[:len(args)-1]

						my_type := self.Class.Name

						other := args[0].(*LigmaInstance)
						other_type := other.Class.Name

						switch my_type {
							case "int":
								switch other_type {
									case "int":
										my_val := self.Fields["value"].(*LigmaInteger).Value
										other_val := other.Fields["value"].(*LigmaInteger).Value
										return nativeBoolToBooleanObject(my_val < other_val)
									case "float":
										my_val := float64(self.Fields["value"].(*LigmaInteger).Value)
										other_val := other.Fields["value"].(*LigmaFloat).Value
										return nativeBoolToBooleanObject(my_val < other_val)
								}
							case "float":
								switch other_type {
									case "int":
										my_val := self.Fields["value"].(*LigmaFloat).Value
										other_val := float64(other.Fields["value"].(*LigmaInteger).Value)
										return nativeBoolToBooleanObject(my_val < other_val)
									case "float":
										my_val := self.Fields["value"].(*LigmaFloat).Value
										other_val := other.Fields["value"].(*LigmaFloat).Value
										return nativeBoolToBooleanObject(my_val < other_val)
								}
						}
						return NewError("unsupported operand type(s) for <: '%s' and '%s'", my_type, other_type)
					},
				},
			},
		},
		Superclasses: []*LigmaClass{builtinsClasses["type"]},
	}


	// add new class "int"
	builtinsClasses["int"] = &LigmaClass{
		Name: "int",
		Methods: ClassMethods{
			BuiltinMethods: map[string]*BuiltinClassMethod{
				"init":{
					Literal: "init",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						args = args[:len(args)-1]
						
						arg_type := args[0].Type()

						switch arg_type {
						case "INTEGER":
							self.Fields["value"] = args[0]
						}
						
						return nil
					},
					NumArgs: 1,
				},
				"__repr__": {
					Literal: "__repr__",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						return &LigmaString{Value: self.Fields["value"].(*LigmaInteger).Inspect()}
					},
				},
				"__str__": {
					Literal: "__str__",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						return &LigmaString{Value: self.Fields["value"].(*LigmaInteger).Inspect()}
					},
				},

			},
			UserDefinedMethods: map[string]*LigmaFunction{},
		},

		Superclasses: []*LigmaClass{builtinsClasses["Number"]},
	}
	
	// float
	builtinsClasses["float"] = &LigmaClass{
		Name: "float",
		Methods: ClassMethods{
			BuiltinMethods: map[string]*BuiltinClassMethod{
				"init":{
					Literal: "init",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						args = args[:len(args)-1]
						
						if len(args) == 0 {
							self.Fields["value"] = &LigmaFloat{Value: 0}
						} else {
							switch arg := args[0].(type) {
							case *LigmaFloat:
								self.Fields["value"] = &LigmaFloat{Value: arg.Value}
							}
						}
						
						return nil
					},
					NumArgs: 1,
				},
				"__repr__": {
					Literal: "__repr__",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						return &LigmaString{Value: self.Fields["value"].(*LigmaFloat).Inspect()}
					},
				},
				"__str__": {
					Literal: "__str__",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						return &LigmaString{Value: self.Fields["value"].(*LigmaFloat).Inspect()}
					},
				},
			},
		},
		Superclasses: []*LigmaClass{builtinsClasses["Number"]},
	}

	// container (for lists, strings, etc)
	builtinsClasses["container"] = &LigmaClass{
		Name: "container",
		Methods: ClassMethods{
			BuiltinMethods: map[string]*BuiltinClassMethod{
				"__repr__": {
					Literal: "__repr__",
					Fn: func(args ...LigmaObject) LigmaObject {
						//self := args[len(args)-1].(*LigmaInstance)
						return &LigmaString{Value: "<container instance>"}
					},
				},
				"__str__": {
					Literal: "__str__",
					Fn: func(args ...LigmaObject) LigmaObject {
						//self := args[len(args)-1].(*LigmaInstance)
						return &LigmaString{Value: "<container instance>"}
					},
				},
				"slice": {
					Literal: "slice",
					Fn: func(args ...LigmaObject) LigmaObject {
						return NewError("Not implemented")
					},
					NumArgs: 2,
				},
				"replace": {
					Literal: "replace",
					Fn: func(args ...LigmaObject) LigmaObject {
						return NewError("Not implemented")
					},
					NumArgs: 2,
				},

				"split": {
					Literal: "split",
					Fn: func(args ...LigmaObject) LigmaObject {
						return NewError("Not implemented")
					},	
					NumArgs: 1,
				},
				"__get__": {
					Literal: "__get__",
					Fn: func(args ...LigmaObject) LigmaObject {
						return NewError("Not implemented")
					},
					NumArgs: 1,
				},
			},
		},
		Superclasses: []*LigmaClass{builtinsClasses["type"]},
	}

	// list
	builtinsClasses["list"] = &LigmaClass{
		Name: "list",
		Methods: ClassMethods{
			BuiltinMethods: map[string]*BuiltinClassMethod{
				"init":{
					Literal: "init",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						args = args[:len(args)-1]

						if len(args) == 0 {
							self.Fields["value"] = &LigmaList{Elements: []LigmaObject{}}
						} else {
							switch arg := args[0].(type) {
							case *LigmaList:
								self.Fields["value"] = &LigmaList{Elements: arg.Elements}
							}
						}

						return nil
					},
					NumArgs: 1,
				},

				"__repr__": {
					Literal: "__repr__",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						elements := self.Fields["value"].(*LigmaList).Elements
						var out []string
						for _, elem := range elements {
							// check if the element has a __repr__ method
							repr, ok := elem.(*LigmaInstance).Get("__repr__")
							if ok {
								repr_callable := repr.(LigmaCallable)
								res := repr_callable.Call(nil, elem)
								out = append(out, res.(*LigmaString).Value)
							}
						}
						return &LigmaString{Value: "[" + strings.Join(out, ", ") + "]"}
					},
				},

				"__get__": {
					Literal: "__get__",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						index := args[0].(*LigmaInstance).Fields["value"].(*LigmaInteger).Value
						elements := self.Fields["value"].(*LigmaList).Elements

						if index < 0 || index >= int64(len(elements)) {
							return NewError("index out of range")
						}

						return elements[index]
					},
					NumArgs: 1,
				},
			},
		},
		Superclasses: []*LigmaClass{builtinsClasses["container"]},
	}

	builtinsClasses["map"] = &LigmaClass{
		Name: "map",
		Methods: ClassMethods{
			BuiltinMethods: map[string]*BuiltinClassMethod{
				"init":{
					Literal: "init",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						args = args[:len(args)-1]

						if len(args) == 0 {
							self.Fields["value"] = &LigmaMap{Pairs: map[MapKey]MapPair{}}
						} else {
							switch arg := args[0].(type) {
							case *LigmaMap:
								self.Fields["value"] = &LigmaMap{Pairs: arg.Pairs}
							}
						}
						return nil
					},
					NumArgs: 1,
				},

				"__repr__": {
					Literal: "__repr__",
					Fn: func(args ...LigmaObject) LigmaObject {
						return NewError("Not implemented")
					},
				},

				"__get__": {
					Literal: "__get__",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						args = args[:len(args)-1]

						mapObj := self.Fields["value"].(*LigmaMap)

						index := args[0]
						key := index.(*LigmaInstance).Fields["value"].(LigmaHashable)

						pair, ok := mapObj.Pairs[key.MapKey()]
						if !ok {
							return NewError("key not found")
						}

						return pair.Value

					},
					NumArgs: 1,
				},
			},
		},
		Superclasses: []*LigmaClass{builtinsClasses["container"]},
	}






	builtinsClasses["str"] = &LigmaClass{
		Name: "str",
		Methods: ClassMethods{
			BuiltinMethods: map[string]*BuiltinClassMethod{
				"init":{
					Literal: "init",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						args = args[:len(args)-1]

						if len(args) == 0 {
							self.Fields["value"] = &LigmaString{Value: ""}
						} else {
							switch arg := args[0].(type) {
							case *LigmaString:
								self.Fields["value"] = &LigmaString{Value: arg.Value}
							}
						}
						return nil
					},
					NumArgs: 1,
				},
				"slice": {
					Literal: "slice",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						args = args[:len(args)-1]

						start := args[0].(*LigmaInstance).Fields["value"].(*LigmaInteger).Value
						end := args[1].(*LigmaInstance).Fields["value"].(*LigmaInteger).Value
						str := self.Fields["value"].(*LigmaString).Value

						if start < 0 || start >= int64(len(str)) || end < 0 || end >= int64(len(str)) {
							return NewError("index out of range")
						}

						return builtinsClasses["str"].Call(nil, &LigmaString{Value: str[start:end]})
					},
					NumArgs: 2,
				},
				"replace": {
					Literal: "replace",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						args = args[:len(args)-1]

						old := args[0].(*LigmaInstance).Fields["value"].(*LigmaString).Value
						new := args[1].(*LigmaInstance).Fields["value"].(*LigmaString).Value
						str := self.Fields["value"].(*LigmaString).Value

						return builtinsClasses["str"].Call(nil, &LigmaString{Value: strings.Replace(str, old, new, -1)})
					},
					NumArgs: 2,
				},
				"split": {
					Literal: "split",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						args = args[:len(args)-1]

						sep := args[0].(*LigmaInstance).Fields["value"].(*LigmaString).Value
						str := self.Fields["value"].(*LigmaString).Value

						var parts []LigmaObject
						for _, part := range strings.Split(str, sep) {
							parts = append(parts, builtinsClasses["str"].Call(nil, &LigmaString{Value: part}))
						}

						return builtinsClasses["list"].Call(nil, parts...)
					},
					NumArgs: 1,
				},
				"__get__": {
					Literal: "__get__",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						index := args[0].(*LigmaInstance).Fields["value"].(*LigmaInteger).Value
						str := self.Fields["value"].(*LigmaString).Value

						if index < 0 || index >= int64(len(str)) {
							return NewError("index out of range")
						}

						return builtinsClasses["str"].Call(nil, &LigmaString{Value: string(str[index])})

					},
					NumArgs: 1,
				},
				"__repr__": {
					Literal: "__repr__",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						str := self.Fields["value"].(*LigmaString)

						return str
					},
				},

				"__str__": {
					Literal: "__str__",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						str := self.Fields["value"].(*LigmaString)

						return str
					},
				},
				"__len__": {
					Literal: "__len__",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						str_list := self.Fields["value"].(*LigmaList)
						return builtinsClasses["int"].Call(nil, &LigmaInteger{Value: int64(len(str_list.Elements))})
					},
				},

				"__add__": {
					Literal: "__add__",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						args = args[:len(args)-1]

						my_str := self.Fields["value"].(*LigmaString).Value
						other := args[0].(*LigmaInstance)
						other_str := other.Fields["value"].(*LigmaString).Value

						return builtinsClasses["str"].Call(nil, &LigmaString{Value: my_str + other_str})
					},
				},

				"__eq__": {
					Literal: "__eq__",
					Fn: func(args ...LigmaObject) LigmaObject {
						self := args[len(args)-1].(*LigmaInstance)
						args = args[:len(args)-1]

						my_str := self.Fields["value"].(*LigmaString).Value
						other := args[0].(*LigmaInstance)
						other_str := other.Fields["value"].(*LigmaString).Value

						return nativeBoolToBooleanObject(my_str == other_str)
					},
				},
				
			},
		},
		Superclasses: []*LigmaClass{builtinsClasses["container"]},
	}
}