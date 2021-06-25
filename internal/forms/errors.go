package forms

type formErrors map[string][]string

// Add adds an error message for the given field
func (fe formErrors) Add(field, message string)  {
	fe[field] = append(fe[field], message)
}

// Get returns the first error message
func (fe formErrors) Get(field string) string {
	es := fe[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
