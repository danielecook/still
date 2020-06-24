package validate

import (
	"errors"
	"fmt"
)

var lastValMap = make(map[string]interface{})

func last(args ...interface{}) (interface{}, error) {
	lastVal := args[0].(string)

	if lastValMap[lastVal] == nil {
		lastValMap[lastVal] = args[1]
		return (bool)(true), nil
	}
	// Get last value
	val := lastValMap[lastVal]
	if val != args[1] {
		// Now set to current value
		lastValMap[lastVal] = args[1]
		err := fmt.Sprintf("%v != %v", args[1], val)
		return (bool)(false), errors.New(err)
	}
	return (bool)(val == args[1]), nil
}
