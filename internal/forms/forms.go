package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

// Form creates a custom form struct, and embeds a Url.value struct
type Form struct {
	url.Values
	Errors errors
}

// New returns a new instance of Form
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// Required checks if argument fields are not empty
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		val := f.Get(field)
		if strings.TrimSpace(val) == "" {
			f.Errors.Add(field, "this field cannot be blank")
		}
	}
}

func (f *Form) Has(field string) bool {
	return f.Get(field) != ""
}

// MinLen checks if provided string is longer than provided len
func (f *Form) MinLen(field string, length int) bool {
	fd := f.Get(field)
	if len(fd) < length {
		msg := fmt.Sprintf("must be more than %d characters long", length)
		f.Errors.Add(field, msg)
		return false
	}
	return true
}

// checks for valid email address
func (f *Form) ValidateEmail() {
	errs := validate.Var(f.Get("email"), "email")
	if errs != nil {
		f.Errors.Add("email", "must be a valid email")
	}
}
