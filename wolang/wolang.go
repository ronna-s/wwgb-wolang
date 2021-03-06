package wolang

import (
	"fmt"
)

func Eval(expr interface{}) (result interface{}, err error) {
	switch v := expr.(type) {
	case string:
		return v, nil
	case int:
		return v, nil
	case []interface{}:
		return evalFCall(v)
	}
	return nil, fmt.Errorf("error: Illegal expression type %T", expr)
}

func evalFCall(expr []interface{}) (interface{}, error) {
	funcName := expr[0].(string)
	arguments := expr[1:]

	if extendedFunctions[funcName] != nil {
		return extendedFunctions[funcName].Call(arguments)
	}

	for ind, arg := range arguments {
		if argarray, ok := arg.([]interface{}); ok {
			nestresult, err := evalFCall(argarray)
			if err != nil {
				fmt.Printf("error: Nested function '%s' evaluation failed", arg)
			}
			arguments[ind] = nestresult
		}
	}

	switch funcName {
	case "+":
		return plus(arguments)
	case "concat":
		return strconcat(arguments)
	}
	return nil, fmt.Errorf("error: Function '%s' is not defined", funcName)
}
