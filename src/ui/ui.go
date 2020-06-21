package ui

import (
	"log"
	"os"

	"github.com/logrusorgru/aurora"
)

// Warning - Print a warning log statement
func Warning(s string) {
	l := log.New(os.Stderr, "", 0)
	l.Println(aurora.Bold(aurora.Yellow("Warning")), s)
}
