package v8

// Field represents a field, its name and
// its format (optional).
type Field struct {
	Field           string
	Format          string
	IncludeUnmapped bool
}

// Source serializes the Field into JSON.
func (d Field) Source() (any, error) {
	if d.Format == "" && d.IncludeUnmapped == false {
		return d.Field, nil
	}

	m := map[string]any{
		"field": d.Field,
	}

	if d.Format != "" {
		m["format"] = d.Format
	}
	if d.IncludeUnmapped == true {
		m["include_unmapped"] = true
	}

	return m, nil
}

// Fields is a slice of Field instances.
type Fields []Field

// Source serializes the Fields into JSON.
func (d Fields) Source() (any, error) {
	if d == nil {
		return nil, nil
	}
	v := make([]any, 0)
	for _, f := range d {
		src, err := f.Source()
		if err != nil {
			return nil, err
		}
		v = append(v, src)
	}
	return v, nil
}
