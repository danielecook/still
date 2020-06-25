package output

import (
	"fmt"
	"strings"

	"github.com/danielecook/still/src/schema"
	"github.com/logrusorgru/aurora"
)

type ValidCol struct {
	Name    string
	IsValid int // 0=not checked; 1=true; 2=false
	NErrs   int
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

// PrintSummary - Prints result summary table
func PrintSummary(validColSet []ValidCol, sch schema.SchemaRules) {
	// Output header
	fmt.Println(strings.Repeat("-", 50))
	fmt.Println(
		aurora.Sprintf(aurora.Bold("%2s\t%-30s\t%10s"),
			"✓",
			"Check",
			"Errors"))
	fmt.Println(strings.Repeat("-", 50))

	if sch.CheckOrdered {
		directiveLine(sch.Ordered, "@ordered")
	}
	if sch.CheckFixed {
		directiveLine(sch.Fixed, "@fixed")
	}

	fmt.Println(strings.Repeat("-", 50))

	nPass := 0
	nChecked := 0
	totalErrs := 0
	var check aurora.Value
	var errs aurora.Value
	for _, col := range validColSet {
		if col.IsValid == 0 {
			check = aurora.Reset("")
			errs = aurora.Reset("")
		} else if col.IsValid == 1 {
			check = aurora.Bold(aurora.Green("✓"))
			errs = aurora.Reset("0")
			nPass++
			nChecked++
		} else if col.IsValid == 2 {
			check = aurora.Bold(aurora.Red("𐄂"))
			errs = aurora.Bold(aurora.Red(col.NErrs))
			nChecked++
			totalErrs += col.NErrs
		}
		fmt.Printf("%2s\t%-30s\t%10v\n",
			check,
			col.Name,
			errs)
	}

	// Summary Line
	fmt.Printf("%s\n", strings.Repeat("-", 50))
	if totalErrs == 0 && sch.Errors == 0 {
		fmt.Printf("%2s\t%-40s\t%10v\n",
			aurora.Bold(aurora.Green("✓")),
			fmt.Sprintf("%-4v",
				aurora.Green("PASS")),
			aurora.Green(totalErrs))
	} else {
		fmt.Printf("%2s\t%-40s\t%10v\n",
			aurora.Bold(aurora.Red("𐄂")),
			fmt.Sprintf("%-4v",
				aurora.Red("FAIL")),
			aurora.Red(totalErrs))
	}
}
