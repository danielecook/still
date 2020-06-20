package validate

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/danielecook/still/src/utils"
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

// strings
func toLower(args ...interface{}) (interface{}, error) {
	return (string)(strings.ToLower(args[0].(string))), nil
}

func toUpper(args ...interface{}) (interface{}, error) {
	return (string)(strings.ToUpper(args[0].(string))), nil
}

// Dates
func parseDate(args ...interface{}) (interface{}, error) {
	val, err := dateparse.ParseAny(fmt.Sprintf("%s", args[0]))
	utils.Check(err)
	return (time.Time)(val), nil
}

func replace(args ...interface{}) (interface{}, error) {
	val := strings.Replace(args[0].(string), args[1].(string), args[2].(string), -1)
	return (string)(val), nil
}

// var lastMap = map[string]map[string]int{}
// func last(args ...interface{}) (interface{}, error) {
// 	// Checks for an element present in a set.

// 	return (bool)(err == nil), nil
// }

// replace

// coalesce
