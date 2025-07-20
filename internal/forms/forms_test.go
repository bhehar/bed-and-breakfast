package forms

import (
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	form := New(url.Values{})
	form.Errors.Add("test_field", "this field cannot be blank")

	isValid := form.Valid()
	if isValid {
		t.Error("form valid when it should be invalid")
	}
}

func TestForm_Required(t *testing.T) {
	form := New(url.Values{})

	form.Required("a", "b")

	if len(form.Errors) != 2 {
		t.Errorf("expected 2 errors but got %d", len(form.Errors))
	}

	form = New(url.Values{})
	form.Add("a", "value-a")
	form.Add("b", "value-b")

	form.Required("a", "b")

	if len(form.Errors) != 0 {
		t.Errorf("expected 0 errors but got %d", len(form.Errors))
	}

}

func TestForm_Has(t *testing.T) {
	// form does not have field
	// r.ParseForm()
	// form := New(r.PostForm)
	form := New(url.Values{})

	if form.Has("key-a") != false {
		t.Error("got true for field that doesn't exist")
	}

	// form does have field
	form = New(url.Values{})
	form.Set("key-a", "val-a")

	if form.Has("key-a") == false {
		t.Error("got false for field that does exist")
	}
}

func TestForm_MinLen(t *testing.T) {
	// len is correct so no error
	form := New(url.Values{})
	form.Set("key-a", "val-a")

	if form.MinLen("key-a", 2) != true {
		t.Errorf("expected true got false")
	}

	if form.Errors.GetFirst("key-a") != "" {
		t.Errorf("expected no error but got one")
	}

	if form.MinLen("key-a", 8) != false {
		t.Errorf("expected false got true")
	}

	if form.Errors.GetFirst("key-a") == "" {
		t.Errorf("expected error but didn't get one")
	}
}

func TestForm_ValidateEmail(t *testing.T) {

	form := New(url.Values{})
	form.Set("email", "valid@email.com")

	form.ValidateEmail()
	if form.Valid() == false {
		t.Errorf("form tested invalid with good email")
	}

	form = New(url.Values{})
	form.Set("email", "not-a-heckin-email")

	form.ValidateEmail()
	if form.Valid() == true {
		t.Errorf("form tested valid with bad email")
	}

}
