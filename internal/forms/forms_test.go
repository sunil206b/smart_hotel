package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/something", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/something", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postData := url.Values{}
	postData.Add("a", "a")
	postData.Add("b", "b")
	postData.Add("c", "c")

	r = httptest.NewRequest("POST", "/something", nil)
	r.PostForm = postData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	postData := url.Values{}
	form := New(postData)
	hasValue := form.Has("a")
	if hasValue {
		t.Error("form has value where it should not")
	}

	postData = url.Values{}
	postData.Add("a", "a")
	form = New(postData)
	hasValue = form.Has("a")
	if !hasValue {
		t.Error("form does not have value where it should have")
	}
}

func TestForm_MinLength(t *testing.T) {
	postData := url.Values{}
	form := New(postData)
	isValid := form.MinLength("a", 5)
	if isValid {
		t.Error("field has min length where it should not")
	}

	err := form.Errors.Get("a")
	if err == "" {
		t.Error("should have an error, but did not get one")
	}

	postData = url.Values{}
	postData.Add("a", "abc")
	form = New(postData)
	isValid = form.MinLength("a", 3)
	if !isValid {
		t.Error("field should have min length, but shows not")
	}
	err = form.Errors.Get("a")
	if err != "" {
		t.Error("should not have an error, but got one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postData := url.Values{}
	form := New(postData)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("form has valid email, where it should not")
	}

	postData = url.Values{}
	postData.Add("email", "a@a.com")
	form = New(postData)
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("form does not have valid email, where it should")
	}
}
