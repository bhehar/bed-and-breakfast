package render

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var templateCache map[string]*template.Template

func RenderTemplateBasic(w http.ResponseWriter, t string) {
	var tmpl *template.Template
	var err error

	// check to see if we already have the template in our cache
	_, inMap := templateCache[t]
	if !inMap {
		// need to read from disk and parse template
		log.Printf("creating %s template and adding to cache.\n", t)
		err = createTemplateCacheBasic(t)
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		// we have template in cache
		log.Printf("using cached %s template\n", t)
	}

	tmpl = templateCache[t]
	err = tmpl.Execute(w, nil)
	if err != nil {
		fmt.Println("error parsing templated", err.Error())
	}
}

func createTemplateCacheBasic(t string) error {
	templates := []string {
		"./templates/" + t,
		"./templates/base.layout.tmpl",
	}
	
	// parse the templates
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}

	// add template to cache
	templateCache[t] = tmpl
	return nil
}
