package output

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/danielecook/still/src/schema"
	"github.com/logrusorgru/aurora"
)

type ValidCol struct {
	Name    string
	IsValid int // 0=not checked; 1=true; 2=false
	NErrs   int
	NNA     int // NA Count
	NEMPTY  int // Empty Count
	NVALID  int
}

func directiveLine(pass bool, name string) {
	if pass {
		fmt.Println(
			aurora.Sprintf("%2s\t%-30s",
				aurora.Bold(aurora.Green("✓")),
				name),
		)
	} else {
		fmt.Println(
			aurora.Sprintf("%2s\t%-30s",
				aurora.Bold(aurora.Red("𐄂")),
				name),
		)
	}
}

func line() {
	fmt.Println(strings.Repeat("-", 82))
}

const lineFormat = "%2s\t%-25s\t%7s\t%7s\t%7v\t%7v\n"

func getField(s ValidCol, field string) int {
	r := reflect.ValueOf(s)
	f := reflect.Indirect(r).FieldByName(field)
	return int(f.Int())
}

func sumField(columns []ValidCol, field string) int {
	s := 0
	for _, col := range columns {
		s += getField(col, field)
	}
	return s
}

// PrintSummary - Prints result summary table
func PrintSummary(validColSet []ValidCol, sch schema.SchemaRules) {
	// Output header
	line()
	fmt.Printf(
		aurora.Sprintf(aurora.Bold(lineFormat),
			"✓",
			"Check",
			"NA",
			"EMPTY",
			"Errors",
			"Valid"))
	line()

	if sch.CheckOrdered {
		directiveLine(sch.Ordered, "@ordered")
	}
	if sch.CheckFixed {
		directiveLine(sch.Fixed, "@fixed")
	}

	line()
	totalErrs := 0
	var check aurora.Value
	var na string
	var empty string
	var valid aurora.Value
	var errs aurora.Value
	for _, col := range validColSet {
		if col.IsValid == 0 {
			check = aurora.Reset("")
			errs = aurora.Reset("")
			na = ""
			empty = ""
			valid = aurora.Reset("")
		} else if col.IsValid >= 1 {
			na = fmt.Sprintf("%d", col.NNA)
			empty = fmt.Sprintf("%d", col.NEMPTY)
			valid = aurora.Green(fmt.Sprintf("%d", col.NVALID))
			if col.IsValid == 1 {
				check = aurora.Bold(aurora.Green("✓"))
				errs = aurora.Green("0")
			} else if col.IsValid == 2 {
				check = aurora.Bold(aurora.Red("𐄂"))
				errs = aurora.Bold(aurora.Red(col.NErrs))
				totalErrs += col.NErrs
			}
		}
		fmt.Printf(lineFormat,
			check,
			col.Name,
			na,
			empty,
			errs,
			valid)
	}

	// Summary Line
	line()
	var status aurora.Value
	na = fmt.Sprintf("%d", sumField(validColSet, "NNA"))
	empty = fmt.Sprintf("%d", sumField(validColSet, "NEMPTY"))
	valid = aurora.Green(fmt.Sprintf("%d", sumField(validColSet, "NVALID")))
	if totalErrs == 0 && sch.Errors == 0 {
		check = aurora.Bold(aurora.Green("✓"))
		status = aurora.Green("PASS")
	} else {
		check = aurora.Bold(aurora.Red("𐄂"))
		status = aurora.Red("FAIL")
	}
	fmt.Printf(lineFormat,
		check,
		status,
		na,
		empty,
		aurora.Red(sumField(validColSet, "NErrs")),
		valid,
	)
}
