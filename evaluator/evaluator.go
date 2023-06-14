package evaluator

import (
	"fmt"
	"github.com/GzArthur/interpreter/ast"
	"github.com/GzArthur/interpreter/object"
)

var (
	TRUE_OBJ  = &object.Boolean{Value: true}
	FALSE_OBJ = &object.Boolean{Value: false}
	NULL_OBJ  = &object.Null{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch n := node.(type) {
	case *ast.Program:
		return evalProgram(n, env)
	case *ast.ExpressionStatement:
		return Eval(n.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}
	case *ast.Boolean:
		return booleanNativeToObj(n.Value)
	case *ast.PrefixExpression:
		rightExpr := Eval(n.RightExpr, env)
		if isErrorObject(rightExpr) {
			return rightExpr
		}
		return evalPrefixExpression(n.Operator, rightExpr)
	case *ast.InfixExpression:
		leftExpr := Eval(n.LeftExpr, env)
		if isErrorObject(leftExpr) {
			return leftExpr
		}
		rightExpr := Eval(n.RightExpr, env)
		if isErrorObject(rightExpr) {
			return rightExpr
		}
		return evalInfixExpression(n.Operator, leftExpr, rightExpr)
	case *ast.BlockStatement:
		return evalBlockStatement(n, env)
	case *ast.IfExpression:
		return evalIfExpression(n, env)
	case *ast.ReturnStatement:
		val := Eval(n.ReturnValue, env)
		if isErrorObject(val) {
			return val
		}
		return &object.Return{Value: val}
	case *ast.LetStatement:
		val := Eval(n.Value, env)
		if isErrorObject(val) {
			return val
		}
		env.Set(n.Name.Value, val)
	case *ast.Identifier:
		return evalIdentifier(n, env)
	case *ast.FunctionLiteral:
		params := n.Parameters
		body := n.Body
		return &object.Function{
			Parameters: params,
			Body:       body,
			Env:        env,
		}
	case *ast.CallExpression:
		function := Eval(n.Function, env)
		if isErrorObject(function) {
			return function
		}
		args := evalExpressions(n.Arguments, env)
		if len(args) == 1 && isErrorObject(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	}
	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var obj object.Object
	for _, stmt := range program.Statements {
		obj = Eval(stmt, env)
		switch res := obj.(type) {
		case *object.Return:
			return res.Value
		case *object.Error:
			return res
		}
	}
	return obj
}

func evalBlockStatement(bStmt *ast.BlockStatement, env *object.Environment) object.Object {
	var obj object.Object
	for _, stmt := range bStmt.Statements {
		obj = Eval(stmt, env)
		if obj != nil {
			objType := obj.Type()
			if objType == object.RETURN_OBJ || objType == object.ERROR_OBJ {
				return obj
			}
		}
	}
	return obj
}

func evalPrefixExpression(op string, expr object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOperationExpression(expr)
	case "-":
		return evalMinusOperationExpression(expr)
	default:
		return newError(fmt.Sprintf("unknown operator: %s%s", op, expr.Type()))
	}
}

func evalInfixExpression(op string, lExpr object.Object, rExpr object.Object) object.Object {
	if lExpr.Type() == object.INTEGER_OBJ && rExpr.Type() == object.INTEGER_OBJ {
		return evalIntegerInfixExpression(op, lExpr, rExpr)
	} else if op == "==" {
		return booleanNativeToObj(lExpr == rExpr)
	} else if op == "!=" {
		return booleanNativeToObj(lExpr != rExpr)
	} else if lExpr.Type() != rExpr.Type() {
		return newError(fmt.Sprintf("type mismatch: %s %s %s", lExpr.Type(), op, rExpr.Type()))
	} else {
		return newError(fmt.Sprintf("unknown operator: %s %s %s", lExpr.Type(), op, rExpr.Type()))
	}
}

func evalIfExpression(node *ast.IfExpression, env *object.Environment) object.Object {
	condObj := Eval(node.Condition, env)
	if isErrorObject(condObj) {
		return condObj
	}
	if isTruth(condObj) {
		return Eval(node.Consequence, env)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, env)
	} else {
		return NULL_OBJ
	}
}

func evalBangOperationExpression(expr object.Object) object.Object {
	switch expr {
	case TRUE_OBJ:
		return FALSE_OBJ
	case FALSE_OBJ, NULL_OBJ:
		return TRUE_OBJ
	default:
		return FALSE_OBJ
	}
}

func evalMinusOperationExpression(expr object.Object) object.Object {
	if expr.Type() != object.INTEGER_OBJ {
		return newError(fmt.Sprintf("unknown operator: -%s", expr.Type()))
	}
	return &object.Integer{Value: -expr.(*object.Integer).Value}
}

func evalIntegerInfixExpression(op string, lExpr object.Object, rExpr object.Object) object.Object {
	lValue := lExpr.(*object.Integer).Value
	rValue := rExpr.(*object.Integer).Value
	switch op {
	case "+":
		return &object.Integer{Value: lValue + rValue}
	case "-":
		return &object.Integer{Value: lValue - rValue}
	case "*":
		return &object.Integer{Value: lValue * rValue}
	case "/":
		return &object.Integer{Value: lValue / rValue}
	case "<":
		return booleanNativeToObj(lValue < rValue)
	case ">":
		return booleanNativeToObj(lValue > rValue)
	case "==":
		return booleanNativeToObj(lValue == rValue)
	case "!=":
		return booleanNativeToObj(lValue != rValue)
	default:
		return newError(fmt.Sprintf("unknown operator: %s %s %s", lExpr.Type(), op, rExpr.Type()))
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError(fmt.Sprintf("identifier not found: %s", node.Value))
	}
	return val
}

func evalExpressions(expr []ast.Expression, env *object.Environment) []object.Object {
	var res []object.Object
	for _, e := range expr {
		evaluated := Eval(e, env)
		if isErrorObject(evaluated) {
			return []object.Object{evaluated}
		}
		res = append(res, evaluated)
	}
	return res
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return newError(fmt.Sprintf("not a function: %s", fn.Type()))
	}
	extendedEnv := extendedFnEnv(function, args)
	evaluated := Eval(function.Body, extendedEnv)
	return unwrapReturnValue(evaluated)
}

func extendedFnEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewWrappedEnv(fn.Env)
	for i, p := range fn.Parameters {
		env.Set(p.Value, args[i])
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.Return); ok {
		return returnValue.Value
	}
	return obj
}

func booleanNativeToObj(input bool) *object.Boolean {
	if input {
		return TRUE_OBJ
	}
	return FALSE_OBJ
}

func isTruth(obj object.Object) bool {
	switch obj {
	case NULL_OBJ, FALSE_OBJ:
		return false
	case TRUE_OBJ:
		return true
	default:
		return true
	}
}

func newError(message string) *object.Error {
	return &object.Error{Message: message}
}

func isErrorObject(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
