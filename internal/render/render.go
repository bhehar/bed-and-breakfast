package render

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/bhehar/bed-and-breakfast/internal/config"
	"github.com/bhehar/bed-and-breakfast/internal/models"
	"github.com/justinas/nosurf"
)

var (
	app *config.AppConfig
	pathToPages   = "./templates/*.page.tmpl"
	pathToLayouts = "./templates/*.layout.tmpl"
)

func NewRenderer(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

// Templates using html/template
func Template(w http.ResponseWriter, r *http.Request, tName string, td *models.TemplateData) error {

	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}
	// get requested template from cache
	t, ok := tc[tName]
	if !ok {
		msg := "can't get template from cache"
		// log.Println(msg)
		return errors.New(msg)
	}

	// to make sure the template can be executed
	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)
	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
		return err
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all of the files named *.page.tmpl from ./templates directory
	pages, err := filepath.Glob(pathToPages)
	if err != nil {
		return myCache, err
	}

	// get all the layout files
	layouts, err := filepath.Glob(pathToLayouts)
	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		// log.Println("CreateCacheTemplate() - file name:", name)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		// if we have layout files, parse them
		if len(layouts) > 0 {
			ts, err = ts.ParseGlob(pathToLayouts)
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
