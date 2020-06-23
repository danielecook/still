package schema

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/danielecook/still/src/utils"
	"gopkg.in/yaml.v2"
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
	YAMLData     map[string]interface{}

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

func parseYaml(yamlData string) map[string]interface{} {
	m := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(yamlData), &m)
	utils.Check(err)
	return m
}

// Parse Schema - Entrypoint
func ParseSchema(schemaFile string) SchemaRules {
	file, err := os.Open(schemaFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var Schema = SchemaRules{}
	var commentOpen = false
	commentEnd, err := regexp.Compile("\\/\\/.*$")
	utils.Check(err)

	scanner := bufio.NewScanner(file)

schema:
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " ")
		line = commentEnd.ReplaceAllString(line, "")

		switch {
		// If an empty line, comment has already been removed.
		case line == "":
			continue
		case strings.HasPrefix(line, "/*"):
			commentOpen = true
			continue
		case strings.HasSuffix(line, "*/"):
			commentOpen = false
			continue
		case commentOpen:
			continue
		case line == "---":
			// If dashline, break and read yaml data
			break schema
		}

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
			default:
				Schema.Separater = '\t'
			}
		}

		// Parse columns
		if strings.HasPrefix(line, "@") == false {
			Schema.Columns = append(Schema.Columns, parseColumn(line))
		}
	}

	var yamlString string
	for scanner.Scan() {
		// Read in yaml data
		yamlString += scanner.Text() + "\n"
	}
	Schema.YAMLData = parseYaml(yamlString)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return Schema
}
