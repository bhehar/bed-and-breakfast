package render

import (
	"net/http"
	"testing"

	"github.com/bhehar/bed-and-breakfast/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123")

	result := AddDefaultData(&td, r)
	if result.Flash != "123" {
		t.Error("failed")
	}
}

func TestTemplate(t *testing.T) {
	pathToPages   = "./../../templates/*.page.tmpl"
	pathToLayouts = "./../../templates/*.layout.tmpl"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww = myWriter{}

	err = Template(&ww, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Errorf("error writing template to browser. details: %v", err)
	}
	err = Template(&ww, r, "non-existent.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("did not get an error for non-existent template")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)
	return r, nil
}

func TestCreateTemplateCache(t *testing.T) {
	pathToPages   = "./../../templates/*.page.tmpl"
	pathToLayouts = "./../../templates/*.layout.tmpl"

	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}
