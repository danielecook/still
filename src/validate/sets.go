package validate

import (
	"crypto/sha1"
	"fmt"
	"strings"
)

func any(args ...interface{}) (interface{}, error) {
	if isNA(args[0]) {
		return (bool)(true), nil
	}
	// Checks for an element present in a set.
	for _, val := range args[1:] {
		subVal, ok := val.([]interface{})
		if ok {
			for _, i := range subVal {
				if args[0] == i {
					return (bool)(true), nil
				}
			}
		}
		if args[0] == val {
			return (bool)(true), nil
		}
	}
	return (bool)(false), nil
}

/*
	unique
*/
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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func isSubsetList(args ...interface{}) (interface{}, error) {
	if isNA(args[0]) {
		return (bool)(true), nil
	}
	testVals := strings.Split(args[0].(string), ",")
	var okVals []string
	switch v := args[1].(type) {
	case string:
		okVals = strings.Split(v, args[2].(string))
	case []interface{}:
		okVals = make([]string, len(v))
		for i, v := range v {
			okVals[i] = v.(string)
		}
	}
	for _, i := range testVals {
		if stringInSlice(i, okVals) == false {
			return (bool)(false), nil
		}
	}
	return (bool)(true), nil
}
