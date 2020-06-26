package validate

import (
	"errors"
	"fmt"
)

var lastValMap = make(map[string]interface{})

func identical(args ...interface{}) (interface{}, error) {
	lastKey := args[0].(string)

	if lastValMap[lastKey] == nil {
		lastValMap[lastKey] = args[1]
		return (bool)(true), nil
	}
	// Get last value
	val := lastValMap[lastKey]
	if val != args[1] {
		// Now set to current value
		lastValMap[lastKey] = args[1]
		err := fmt.Sprintf("%v != %v", val, args[1])
		return (bool)(false), errors.New(err)
	}
	return (bool)(val == args[1]), nil
}

func last(args ...interface{}) (interface{}, error) {
	lastKey := args[0].(string)
	if lastValMap[lastKey] == nil {
		lastValMap[lastKey] = args[1]
		return (string)(""), nil
	}
	lastValMap[lastKey] = args[1]
	return lastValMap[lastKey], nil
}
