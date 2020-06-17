package validate

import "math"

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
