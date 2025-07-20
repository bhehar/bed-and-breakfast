package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	mh := myHandler{}
	h := NoSurf(&mh)

	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf("expected type http.Handler got %v", v)
	}
}

func TestSessionLoad(t *testing.T) {
	mh := myHandler{}
	h := SessionLoad(&mh)

	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf("expected type http.Handler got %v", v)
	}
}