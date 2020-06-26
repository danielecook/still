package validate

import (
	"fmt"
	"log"
	"reflect"
)

func is(args ...interface{}) (interface{}, error) {
	if m, _ := isMissing(args[0]); m.(bool) {
		return (bool)(true), nil
	}
	return (bool)(args[0] == args[1]), nil
}

func not(args ...interface{}) (interface{}, error) {
	if m, _ := isMissing(args[0]); m.(bool) {
		return (bool)(true), nil
	}
	return (bool)(args[0] != args[1]), nil
}

func ifElse(args ...interface{}) (interface{}, error) {
	// Assert that both results are the same type
	if reflect.TypeOf(args[1]) != reflect.TypeOf(args[2]) {
		errMsg := fmt.Sprintf("if_else  must return the same type for TRUE and FALSE; Returning %T != %T", args[1], args[2])
		log.Fatal(errMsg)
	}

	if args[0].(bool) {
		return (interface{})(args[1].(interface{})), nil
	}
	return (interface{})(args[2].(interface{})), nil
}
