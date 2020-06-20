package validate

import (
	"math"
)

func strLen(args ...interface{}) (interface{}, error) {
	length := len(args[0].(string))
	return (float64)(length), nil
}

func maxFunc(args ...interface{}) (interface{}, error) {
	var max float64
	for _, n := range args {
		max = math.Max(n.(float64), max)
	}
	return (float64)(max), nil
}

func minFunc(args ...interface{}) (interface{}, error) {
	var min float64
	for _, n := range args {
		min = math.Min(n.(float64), min)
	}
	return (float64)(min), nil
}

// Returns a count of the current arguments
// Uses the uniqueMap as part of sets
func countFunc(args ...interface{}) (interface{}, error) {
	uniqueGroup := args[0].(string)
	digest := digestArgs(args[1:])
	if uniqueMap[uniqueGroup] == nil {
		uniqueMap[uniqueGroup] = make(map[string]int)
	}
	if _, ok := uniqueMap[uniqueGroup][digest]; ok == false {
		// The value has not been seen
		uniqueMap[uniqueGroup][digest] = 0
	}
	uniqueMap[uniqueGroup][digest]++
	return (float64)(uniqueMap[uniqueGroup][digest]), nil
}
