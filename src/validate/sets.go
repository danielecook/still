package validate

import (
	"crypto/sha1"
	"fmt"
)

func any(args ...interface{}) (interface{}, error) {
	// Checks for an element present in a set.
	for _, val := range args[1:] {
		if args[0] == val {
			return (bool)(true), nil
		}
	}
	return (bool)(false), nil
}

var uniqueMap = map[string]map[string]int{}

func digestArgs(args ...interface{}) string {
	h := sha1.New()
	return string(h.Sum([]byte(fmt.Sprintf("%v", args))))
}

func uniqueFunc(args ...interface{}) (interface{}, error) {
	uniqueGroup := args[0].(string)
	digest := digestArgs(args[1:])
	if uniqueMap[uniqueGroup] == nil {
		uniqueMap[uniqueGroup] = make(map[string]int)
	}
	if _, ok := uniqueMap[uniqueGroup][digest]; ok == false {
		// The value has not been seen
		uniqueMap[uniqueGroup][digest] = 1
		return (bool)(true), nil
	}
	uniqueMap[uniqueGroup][digest]++
	return (bool)(false), nil
}

// is_subsetted_list()
