package validate

import (
	"fmt"
	"strconv"
	"strings"
)

func isNumeric(args ...interface{}) (interface{}, error) {
	switch args[0].(type) {
	case float64, float32, int:
		return (bool)(true), nil
	}
	return (bool)(false), nil
}

func isInt(args ...interface{}) (interface{}, error) {
	if _, err := strconv.Atoi(fmt.Sprintf("%s", args[0])); err == nil {
		return (bool)(true), nil
	}
	return (bool)(false), nil
}

func isBool(args ...interface{}) (interface{}, error) {
	val := fmt.Sprintf("%s", args[0])
	BooleanValues := []string{"true", "false"}
	for _, x := range BooleanValues {
		if strings.ToLower(val) == x {
			return (bool)(true), nil
		}
	}
	return (bool)(false), nil
}