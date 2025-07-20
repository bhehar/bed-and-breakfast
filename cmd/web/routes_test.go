package main

import (
	"testing"

	"github.com/bhehar/bed-and-breakfast/internal/config"
	"github.com/go-chi/chi"
)

func TestRoutes(t *testing.T) {

	app := config.AppConfig{}

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// test passes
	default:
		t.Errorf("expected type *chi.Mux got %v", v)
	}
}
