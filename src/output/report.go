package output

import (
	"fmt"
	"strings"

	"github.com/logrusorgru/aurora"
)

type ValidCol struct {
	Name    string
	IsValid bool
	NErrs   int
}

// PrintSummary - Prints result summary table
func PrintSummary(validColSet []ValidCol) {
	fmt.Println(strings.Repeat("-", 50))
	fmt.Println(
		aurora.Sprintf(aurora.Bold("%2s\t%-30s\t%10s"),
			"✓",
			"Column",
			"Errors"))
	fmt.Println(strings.Repeat("-", 50))
	allPass := true
	nPass := 0
	totalErrs := 0
	var check aurora.Value
	var errs aurora.Value
	for _, col := range validColSet {
		if col.IsValid {
			check = aurora.Bold(aurora.Green("✓"))
			errs = aurora.Reset("0")
			nPass++
		} else {
			check = aurora.Bold(aurora.Red("𐄂"))
			allPass = false
			errs = aurora.Bold(aurora.Red(col.NErrs))
			totalErrs += col.NErrs
		}
		fmt.Printf("%2s\t%-30s\t%10v\n",
			check,
			col.Name,
			errs)
	}

	// Summary Line
	fmt.Printf("%s\n", strings.Repeat("-", 50))
	if allPass {
		fmt.Printf("%2s\t%-40s\t%10v\n",
			aurora.Bold(aurora.Green("✓")),
			fmt.Sprintf("%-4v (%d/%d)",
				aurora.Green("PASS"), nPass, len(validColSet)),
			aurora.Green(totalErrs))
	} else {
		fmt.Printf("%2s\t%-40s\t%10v\n",
			aurora.Bold(aurora.Red("𐄂")),
			fmt.Sprintf("%-4v (%d/%d)",
				aurora.Red("FAIL"), nPass, len(validColSet)),
			aurora.Red(totalErrs))
	}
}
