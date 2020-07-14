package validate

/*
	Two pass functions

	(1) First pass through file computes something
	(2) Second pass uses computed values to evaluate an expression
*/

import (
	"fmt"
)

var groupMap = map[string]map[string]int{}

func groupCountFunc(args ...interface{}) (interface{}, error) {
	/*
		(1) hash - stores data for function in specific hash
		(2) group - group column
		(3) count_column - value to count within group
		(4) eqVal - Value to count
	*/
	if m, _ := isMissing(args); m.(bool) {
		return (bool)(true), nil
	}
	hashCol := args[0].(string) // groupHash
	groupCol := fmt.Sprintf("%v", args[1])
	countCol := fmt.Sprintf("%v", args[2])
	eqVal := fmt.Sprintf("%v", args[3])
	if groupMap[hashCol] == nil {
		groupMap[hashCol] = make(map[string]int)
	}
	if _, ok := groupMap[hashCol][groupCol]; ok == false {
		// The value has not been seen
		groupMap[hashCol][groupCol] = 0
	}
	if countCol == eqVal {
		groupMap[hashCol][groupCol]++
	}
	return (map[string]map[string]int)(groupMap), nil
}

func groupCountFuncEval(args ...interface{}) (interface{}, error) {
	if m, _ := isMissing(args); m.(bool) {
		return (bool)(true), nil
	}
	hashCol := args[0].(string) // groupHash
	groupCol := fmt.Sprintf("%v", args[1])
	//warning := errors.New(fmt.Sprintf("Group Count = %d", groupMap[hashCol][groupCol]))
	return (float64)(groupMap[hashCol][groupCol]), nil //warning
}
