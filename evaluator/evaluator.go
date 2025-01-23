package evaluator

import (
	"github.com/ekediala/jian/ast"
	"github.com/ekediala/jian/object"
	"github.com/ekediala/jian/token"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch val := node.(type) {
	case *ast.Program:
		return evalProgram(val.Statements, env)

	case *ast.ExpressionStatement:
		return Eval(val.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(val, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: val.Value}

	case *ast.StringLiteral:
		return &object.String{Value: val.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(val.Value)

	case *ast.PrefixExpression:
		right := Eval(val.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(val.Operator, right)

	case *ast.InfixExpression:
		left := Eval(val.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(val.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(left, val.Operator, right)

	case *ast.IfExpression:
		return evalIfExpression(val, env)

	case *ast.ReturnStatement:
		v := Eval(val.ReturnValue, env)
		if isError(v) {
			return v
		}
		return &object.ReturnValue{Value: v}

	case *ast.LetStatement:
		v := Eval(val.Value, env)
		if isError(v) {
			return v
		}
		env.Set(val.Name.Value, v)

	case *ast.Identifier:
		return evalIdentifier(val, env)

	case *ast.FunctionLiteral:
		return object.NewFunction(val.Parameters, val.Body, env)

	case *ast.CallExpression:
		{
			fn := Eval(val.Function, env)
			if isError(fn) {
				return fn
			}

			if function, ok := fn.(*object.Function); ok {
				if exp, got := len(function.Parameters), len(val.Arguments); exp != got {
					return object.NewError("invalid argument length; expected %d arguments, got %d", exp, got)
				}
			}

			args := evalExpressions(val.Arguments, env)
			if len(args) == 1 && isError(args[0]) {
				return args[0]
			}
			return applyFunction(fn, args)
		}

	case *ast.ArrayLiteral:
		elements := evalExpressions(val.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.IndexExpression:
		{
			left := Eval(val.Left, env)
			if isError(left) {
				return left
			}
			
			index := Eval(val.Index, env)
			if isError(index) {
				return index
			}
			
			return evalIndexOperation(left, index)
		}
	}

	return nil
}

func evalIndexOperation(left object.Object, index object.Object) object.Object {
	switch obj := left.(type) {
	case *object.Array:
		{
			if i, ok := index.(*object.Integer); ok {
				return evalArrayIndexOperation(obj, i)
			}
			return object.NewError("expected index to be *object.Integer, got %T", index)
		}

	default:
		return object.NewError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexOperation(arr *object.Array, index *object.Integer) object.Object {
	if index.Value < 0 || index.Value >= int64(len(arr.Elements)) {
		return NULL
	}

	return arr.Elements[index.Value]
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var r object.Object
	for _, stmt := range stmts {
		r = Eval(stmt, env)
		switch result := r.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return r
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var r object.Object

	for _, stmt := range block.Statements {
		r = Eval(stmt, env)
		if r.Type() == object.RETURN_VALUE || r.Type() == object.ERROR {
			return r
		}
	}

	return r
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case token.BANG:
		return evalBangOperatorExpression(right)
	case token.MINUS:
		return evalMinusOperatorExpression(right)
	default:
		return object.NewError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	if obj, ok := right.(*object.Integer); ok {
		return &object.Integer{Value: -obj.Value}
	}
	return object.NewError("unknown operator: -%s", right.Type())
}

func evalInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	switch {
	case left.Type() != right.Type():
		return object.NewError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	case left.Type() == object.INTEGER:
		return evalIntegerInfixOperation(left.(*object.Integer), operator, right.(*object.Integer))
	case left.Type() == object.STRING:
		return evalStringInfixOperation(left.(*object.String), operator, right.(*object.String))
	case operator == token.EQ:
		return nativeBoolToBooleanObject(left == right)
	case operator == token.NOT_EQ:
		return nativeBoolToBooleanObject(left != right)
	default:
		return object.NewError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalStringInfixOperation(left *object.String, operator string, right *object.String) object.Object {
	switch operator {
	case token.PLUS:
		return &object.String{Value: left.Value + right.Value}
	case token.EQ:
		return nativeBoolToBooleanObject(left.Value == right.Value)
	case token.NOT_EQ:
		return nativeBoolToBooleanObject(left.Value != right.Value)
	default:
		return object.NewError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixOperation(left *object.Integer, operator string, right *object.Integer) object.Object {
	switch operator {
	case token.MINUS:
		return &object.Integer{Value: left.Value - right.Value}
	case token.PLUS:
		return &object.Integer{Value: left.Value + right.Value}
	case token.SLASH:
		return &object.Integer{Value: left.Value / right.Value}
	case token.ASTERISK:
		return &object.Integer{Value: left.Value * right.Value}
	case token.LT:
		return nativeBoolToBooleanObject(left.Value < right.Value)
	case token.GT:
		return nativeBoolToBooleanObject(left.Value > right.Value)
	case token.EQ:
		return nativeBoolToBooleanObject(left.Value == right.Value)
	case token.NOT_EQ:
		return nativeBoolToBooleanObject(left.Value != right.Value)
	default:
		return object.NewError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIfExpression(exp *ast.IfExpression, env *object.Environment) object.Object {
	cond := Eval(exp.Condition, env)
	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return Eval(exp.Consequence, env)
	}

	if exp.Alternative != nil {
		return Eval(exp.Alternative, env)
	}

	return NULL
}

func isTruthy(cond object.Object) bool {
	switch cond {
	case TRUE:
		return true
	case FALSE:
		return false
	case NULL:
		return false
	default:
		return true
	}
}

func isError(obj object.Object) bool {
	return obj != nil && obj.Type() == object.ERROR
}

func evalIdentifier(ident *ast.Identifier, env *object.Environment) object.Object {
	if obj, ok := env.Get(ident.Value); ok {
		return obj
	}

	if obj, ok := builtins[ident.Value]; ok {
		return obj
	}

	return object.NewError("identifier not found: %s", ident.Value)
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var r []object.Object

	for _, exp := range exps {
		val := Eval(exp, env)
		if isError(val) {
			return []object.Object{val}
		}

		r = append(r, val)
	}

	return r
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch obj := fn.(type) {
	case *object.Function:
		{
			env := extendFunctionEnv(obj, args)
			val := Eval(obj.Body, env)
			return unwrapReturnValue(val)
		}
	case *object.Builtin:
		{
			return obj.Fn(args...)
		}
	default:
		return object.NewError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for i, param := range fn.Parameters {
		env.Set(param.Value, args[i])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if rv, ok := obj.(*object.ReturnValue); ok {
		return rv.Value
	}

	return obj
}
