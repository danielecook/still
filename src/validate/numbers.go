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
