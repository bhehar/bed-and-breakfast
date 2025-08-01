package forms

type errors map[string][]string

// Add adds an error message to errors for a give form field
func (e errors) Add(field, msg string) {
	e[field] = append(e[field], msg)
}

// Get returns the first error message
func (e errors) GetFirst(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}