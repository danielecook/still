package validate

import (
	"fmt"
	"strconv"
	"strings"
)

func isNumeric(args ...interface{}) (interface{}, error) {
	if m, _ := isMissing(args[0]); m.(bool) {
		return (bool)(true), nil
	}
	switch args[0].(type) {
	case float64, float32, int:
		return (bool)(true), nil
	}
	return (bool)(false), nil
}

func isInt(args ...interface{}) (interface{}, error) {
	if m, _ := isMissing(args[0]); m.(bool) {
		return (bool)(true), nil
	}
	conv, ok := args[0].(float64)
	if ok {
		if conv == float64(int64(conv)) {
			return (bool)(true), nil
		}
	}
	if _, err := strconv.Atoi(fmt.Sprintf("%s", args[0])); err == nil {
		return (bool)(true), nil
	}
	return (bool)(false), nil
}

func isBool(args ...interface{}) (interface{}, error) {
	if m, _ := isMissing(args[0]); m.(bool) {
		return (bool)(true), nil
	}
	val := fmt.Sprintf("%v", args[0])
	BooleanValues := []string{"true", "false"}
	for _, x := range BooleanValues {
		if strings.ToLower(val) == x {
			return (bool)(true), nil
		}
	}
	return (bool)(false), nil
}

func isString(args ...interface{}) (interface{}, error) {
	if m, _ := isMissing(args[0]); m.(bool) {
		return (bool)(true), nil
	}
	tNum, _ := isNumeric(args[0])
	tInt, _ := isInt(args[0])
	tBool, _ := isBool(args[0])
	if tNum == false &&
		tInt == false &&
		tBool == false {
		return (bool)(true), nil
	}
	return (bool)(false), nil
}
