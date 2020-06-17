package validate

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

func regex(args ...interface{}) (interface{}, error) {
	isMatch, err := regexp.MatchString(args[1].(string), args[0].(string))
	if err != nil {
		log.Fatal(err)
	}
	return (bool)(isMatch), nil
}

func uppercase(args ...interface{}) (interface{}, error) {
	val := args[0].(string)
	isUpper := strings.ToUpper(val) == val
	return (bool)(isUpper), nil
}

func lowercase(args ...interface{}) (interface{}, error) {
	val := args[0].(string)
	isUpper := strings.ToLower(val) == val
	return (bool)(isUpper), nil
}

func starts(args ...interface{}) (interface{}, error) {
	startsWith := strings.HasPrefix(args[0].(string), args[1].(string))
	return (bool)(startsWith), nil
}

func ends(args ...interface{}) (interface{}, error) {
	endsWith := strings.HasSuffix(args[0].(string), args[1].(string))
	return (bool)(endsWith), nil
}

func contains(args ...interface{}) (interface{}, error) {
	isIn := strings.Contains(fmt.Sprintf("%s", args[0]), args[1].(string))
	return (bool)(isIn), nil
}

func length(args ...interface{}) (interface{}, error) {
	s := len(fmt.Sprintf("%s", args[0]))
	minLen := int(args[1].(float64))
	var withinLen bool
	fmt.Println(s)
	if len(args) == 2 {
		withinLen = minLen == s
	} else if len(args) == 3 {
		switch maxLen := args[2].(type) {
		case string:
			if maxLen == "*" {
				withinLen = minLen <= s
			} else {
				log.Fatal("Invalid parameter")
			}
		case float64:
			withinLen = minLen <= s && s <= int(maxLen)
		}
	} else {
		log.Fatal("Length only takes 2 arguments")
	}
	return (bool)(withinLen), nil
}
