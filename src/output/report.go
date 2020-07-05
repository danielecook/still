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

func getField(s schema.Col, field string) int {
	r := reflect.ValueOf(s)
	f := reflect.Indirect(r).FieldByName(field)
	return int(f.Int())
}

func sumField(columns []schema.Col, field string) int {
	s := 0
	for _, col := range columns {
		s += getField(col, field)
	}
	return s
}

// PrintSummary - Prints result summary table
func PrintSummary(colnames []string, sch schema.SchemaRules) {
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

	// Allow for sorting by DATA or SCHEMA
	var colOrder []string
	// Restructure columns into hash map
	ColSet := make(map[string]schema.Col)
	for _, col := range sch.Columns {
		ColSet[col.Name] = col
	}

	if sch.OutputOrder == "schema" {
		colOrder = make([]string, len(sch.Columns))
		// Order columns by schema
		for idx, col := range sch.Columns {
			colOrder[idx] = col.Name
		}
		for _, cn := range colnames {
			// add unchecked columns to the end
			if _, ok := ColSet[cn]; ok == false {
				colOrder = append(colOrder, cn)
			}
		}
	} else {
		// Order by data
		colOrder = make([]string, len(colnames))
		// If sorting by data, add missing columns in schema to end of colnames
		colOrder = colnames
		for _, col := range sch.Columns {
			if col.Status == 3 {
				colOrder = append(colOrder, col.Name)
			}
		}
	}

	line()
	totalErrs := 0
	var check aurora.Value
	var na string
	var empty string
	var colName aurora.Value
	var valid aurora.Value
	var errs aurora.Value

	for _, cname := range colOrder {
		col := ColSet[cname]
		// Add information for columns missing in schema
		colName = aurora.Reset(col.Name)
		if col.Name == "" {
			col.Name = cname
			colName = aurora.Gray(15, col.Name)
		}
		if col.Status == 0 {
			// Not Checked
			check = aurora.Reset("")
			errs = aurora.Reset("")
			na = ""
			empty = ""
			valid = aurora.Reset("")
		} else if col.Status >= 1 {
			na = fmt.Sprintf("%d", col.NNA)
			empty = fmt.Sprintf("%d", col.NEMPTY)
			valid = aurora.Green(fmt.Sprintf("%d", col.NVALID))
			if col.Status == 1 {
				// Valid
				check = aurora.Bold(aurora.Green("✓"))
				errs = aurora.Green("0")
			} else if col.Status == 2 {
				// Invalid
				check = aurora.Bold(aurora.Red("𐄂"))
				errs = aurora.Bold(aurora.Red(col.NErrs))
				totalErrs += col.NErrs
			} else if col.Status == 3 {
				// Missing
				colName = aurora.Gray(15, fmt.Sprintf("%s [missing in data]", col.Name))
				check = aurora.Bold(aurora.Red("𐄂"))
				errs = aurora.Reset("")
				na = ""
				empty = ""
				valid = aurora.Reset("")
				totalErrs += col.NErrs
			}
		}
		fmt.Printf(lineFormat,
			check,
			colName,
			na,
			empty,
			errs,
			valid)
	}

	// Summary Line
	line()
	var status aurora.Value
	na = fmt.Sprintf("%d", sumField(sch.Columns, "NNA"))
	empty = fmt.Sprintf("%d", sumField(sch.Columns, "NEMPTY"))
	valid = aurora.Green(fmt.Sprintf("%d", sumField(sch.Columns, "NVALID")))
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
		aurora.Red(sumField(sch.Columns, "NErrs")),
		valid,
	)
}
