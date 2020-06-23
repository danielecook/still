package validate

func isPositive(args ...interface{}) (interface{}, error) {
	// Checks for an element present in a set.
	if val, ok := args[0].(float64); ok {
		return (bool)(val > 0), nil
	}
	return (bool)(false), nil
}

func isNegative(args ...interface{}) (interface{}, error) {
	// Checks for an element present in a set.
	if val, ok := args[0].(float64); ok {
		return (bool)(val < 0), nil
	}
	return (bool)(false), nil
}

func rangeFunc(args ...interface{}) (interface{}, error) {
	if isNA(args[0]) {
		return (bool)(true), nil
	}
	var between bool
	switch val := args[0].(type) {
	case nil:
		return (bool)(true), nil
	case float64:
		between = args[1].(float64) <= val && val <= args[2].(float64)
	}
	val := args[0].(float64)
	between = args[1].(float64) <= val && val <= args[2].(float64)
	return (bool)(between), nil
}
