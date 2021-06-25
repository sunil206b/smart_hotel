package forms

import (
	"net/http"
	"net/url"
	"strings"
)

// Form creates a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors formErrors
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		formErrors{},
	}
}

// Has checks if form field is in post and not empty
func (f *Form) Has(field string, r *http.Request) bool {
	s := r.Form.Get(field)
	if s == "" {
		f.Errors.Add(field, "This field cannot be blank")
		return false
	}
	return true
}

func (f *Form) Required(fields ...string)  {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field is required")
		}
	}
}

//Valid returns true if there are no errors
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}