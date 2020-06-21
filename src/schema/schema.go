package schema

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Col struct {
	Name        string
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

func parseDirectiveValue(line string) string {
	i := strings.Fields(line)
	if len(i) == 1 {
		log.Fatal(fmt.Sprintf("Missing value for directive %s", i[0]))
	}
	return i[1]
}

func parseDirectiveInt(line string) int {
	i, err := strconv.Atoi(parseDirectiveValue(line))
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func parseDirectiveStrArray(line string) []string {
	fields := strings.Fields(line)
	fields = fields[1:]
	// Set empty strings
	for idx, val := range fields {
		if val == "\"\"" || val == "''" {
			fields[idx] = ""
		}
	}
	return fields
}

func parseColumn(line string) Col {
	/*
		Parses a column specification

		latitude: in_range(-90, 90) # The latitude
	*/
	reColumn, err := regexp.Compile("^([^: ]+):([^$#]+)(#.*$)?")
	if err != nil {
		log.Fatal(err)
	}
	colmatch := reColumn.FindStringSubmatch(line)
	col := Col{}
	col.Name = colmatch[1]
	col.Rule = colmatch[2]
	col.Description = strings.Trim(colmatch[3], "# ")
	return col
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
			case strings.HasPrefix(line, "@na_values"):
				Schema.NA = parseDirectiveStrArray(line)
			case strings.HasPrefix(line, "@"):
				log.Fatal(fmt.Sprintf("%s is an unknown directive", line))
			case strings.HasPrefix(line, "#"):
				continue
			default:
				Schema.Separater = '\t'
			}
		}

		// Parse columns
		if strings.HasPrefix(line, "@") == false {
			Schema.Columns = append(Schema.Columns, parseColumn(line))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return Schema
}
