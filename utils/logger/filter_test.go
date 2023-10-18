package logger

import (
	"regexp"
	"testing"
)

func TestFilterTypeToFunc_Regexp(t *testing.T) {
	r := regexp.MustCompile("^test")
	filterFunc := NewFilter(r).Done()
	if !filterFunc("test123") {
		t.Error("expected true when input string starts with 'test'")
	}
	if filterFunc("123test") {
		t.Error("expected false when input string does not start with 'test'")
	}
}

func TestFilterTypeToFunc_String(t *testing.T) {
	filterFunc := NewFilter("test").Done()
	if !filterFunc("test") {
		t.Error("expected true when input string is exactly 'test'")
	}
	if filterFunc("123test") {
		t.Error("expected false when input string is not exactly 'test'")
	}
}

func TestGenericFilter_Not(t *testing.T) {
	r := regexp.MustCompile("^test")
	gf := NewFilter(r)
	gfNot := gf.Not()

	if gfNot.filter("test123") {
		t.Error("expected false when input string starts with 'test' after Not operation")
	}
	if !gfNot.filter("123test") {
		t.Error("expected true when input string does not start with 'test' after Not operation")
	}
}

func TestFilterTypeToFunc_FilterNone(t *testing.T) {
	filterFunc := FilterMatchNone
	if filterFunc("test") {
		t.Error("expected false for any input when FilterMatchNone is used")
	}
	if filterFunc("123") {
		t.Error("expected false for any input when FilterMatchNone is used")
	}
}

func TestFilterTypeToFunc_FilterAll(t *testing.T) {
	filterFunc := FilterMatchAll
	if !filterFunc("test") {
		t.Error("expected true for any input when FilterMatchAll is used")
	}
	if !filterFunc("123") {
		t.Error("expected true for any input when FilterMatchAll is used")
	}
}

func TestGenericFilter_And(t *testing.T) {
	r1 := regexp.MustCompile("^test")
	r2 := regexp.MustCompile("123$")
	gf := NewFilter(r1)
	gfAnd := gf.And(r2)

	if !gfAnd.filter("test123") {
		t.Error("expected true when input string starts with 'test' and ends with '123' after And operation")
	}
	if gfAnd.filter("123test") {
		t.Error("expected false when input string does not start with 'test' but ends with '123' after And operation")
	}
	if gfAnd.filter("test124") {
		t.Error("expected false when input string starts with 'test' but does not end with '123' after And operation")
	}
}

func TestGenericFilter_Or(t *testing.T) {
	r1 := regexp.MustCompile("^test")
	r2 := regexp.MustCompile("123$")
	gf := NewFilter(r1)
	gfOr := gf.Or(r2).Done()

	if !gfOr("test123") {
		t.Error("expected true when input string starts with 'test' and ends with '123' after Or operation")
	}
	if !gfOr("sobeit123") {
		t.Error("expected true when input string does not start with 'test' but ends with '123' after Or operation")
	}
	if !gfOr("test124") {
		t.Error("expected true when input string starts with 'test' but does not end with '123' after Or operation")
	}
	if !gfOr("test") {
		t.Error("expected true when input string starts with 'test' but does not end with '123' after Or operation")
	}
	if gfOr("124") {
		t.Error("expected false when input string neither starts with 'test' nor ends with '123' after Or operation")
	}
}
