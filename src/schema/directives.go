package schema

func indexOf(name string, data []string) int {
	for k, v := range data {
		if name == v {
			return k
		}
	}
	return -1
}

// IsOrdered - test whether column names are in the
// same order as the schema.
func (s *SchemaRules) IsOrdered(colnames []string) {
	if s.CheckOrdered {
		s.Ordered = true
		idxMax := 0
		for _, x := range s.Columns {
			i := indexOf(x.Name, colnames)
			if i > idxMax {
				idxMax = i
			} else {
				s.Ordered = false
			}
		}
	}
}

// IsFixed - test whether column names match exactly
func (s *SchemaRules) IsFixed(colnames []string) {
	if s.CheckFixed {
		s.Fixed = true
		if len(s.Columns) != len(colnames) {
			s.Fixed = false
			return
		}
		for i := range colnames {
			if colnames[i] != s.Columns[i].Name {
				s.Fixed = false
			}
		}
	}
}
