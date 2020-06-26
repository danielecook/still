package validate

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/araddon/dateparse"
)

func isDate(args ...interface{}) (interface{}, error) {
	if m, _ := isMissing(args[0]); m.(bool) {
		return (bool)(true), nil
	}
	// Checks for an element present in a set.
	_, err := dateparse.ParseStrict(fmt.Sprintf("%s", args[0]))
	return (bool)(err == nil), nil
}

func isDateRelaxed(args ...interface{}) (interface{}, error) {
	if m, _ := isMissing(args[0]); m.(bool) {
		return (bool)(true), nil
	}
	// Checks for an element present in a set.
	_, err := dateparse.ParseAny(fmt.Sprintf("%s", args[0]))
	return (bool)(err == nil), nil
}

func isDateFormat(args ...interface{}) (interface{}, error) {
	if m, _ := isMissing(args[0]); m.(bool) {
		return (bool)(true), nil
	}
	var format = strings.Replace(args[1].(string), "\\", "", -1)
	format = strings.Trim(format, "[]")
	layout, err := dateparse.ParseFormat(format)
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to parse '%s'", args[1]))
	}
	_, err = time.Parse(layout, args[0].(string))
	return (bool)(err == nil), nil
}
