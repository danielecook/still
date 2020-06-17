package schema

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Col struct {
	ColName     string
	Rule        string
	Description string
}

type SchemaRules struct {
	// Directives
	Separater    rune
	ExactColumns int
	Comment      rune
	NA           []string

	// Columns
	Columns []Col
}

func parseDirectiveValue(directive string) string {
	i := strings.Fields(directive)
	if len(i) == 1 {
		log.Fatal(fmt.Sprintf("Missing value for directive %s", i[0]))
	}
	return i[1]
}

func parseDirectiveInt(directive string) int {
	i, err := strconv.Atoi(parseDirectiveValue(directive))
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func parseDirectiveStrArray(directive string) []string {
	return strings.Split(parseDirectiveValue(directive), ",")
}

func setSeparator(sep string) rune {
	switch sep {
	case "TAB":
		return '\t'
	case "\t":
		return '\t'
	}
	return rune(sep[0])
}

// Parse Schema - Entrypoint
func ParseSchema(schemaFile string) SchemaRules {
	file, err := os.Open(schemaFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var Schema = SchemaRules{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Parse Directives
		if strings.HasPrefix(line, "@") {
			switch {
			case strings.HasPrefix(line, "@columns"):
				Schema.ExactColumns = parseDirectiveInt(line)
			case strings.HasPrefix(line, "@sep"):
				Schema.Separater = setSeparator(parseDirectiveValue(line))

			case strings.HasPrefix(line, "@"):
				log.Fatal(fmt.Sprintf("%s is an unknown directive", line))
			default:
				Schema.Separater = '\t'
			}
		}

		// Parse column rules
		if strings.HasPrefix(line, "@") == false {
			Schema.Colnames
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return Schema
}
