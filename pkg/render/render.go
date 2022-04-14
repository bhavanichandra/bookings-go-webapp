package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/bhavanichandra/bookings/pkg/config"
	"github.com/bhavanichandra/bookings/pkg/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

//RenderTemplate renders template using html templates
func RenderTemplate(w http.ResponseWriter, tmpl string, data *models.TemplateData) {

	// Get the template cache from app config
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}
	myTemplate, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from Template Cache")
	}
	buf := new(bytes.Buffer)
	data = AddDefaultData(data)
	_ = myTemplate.Execute(buf, data)
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}
}

//CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	pages, err := filepath.Glob("./templates/*.page.tmpl") // Get all the filenames end with .page.tmpl
	if err != nil {
		return myCache, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
			myCache[name] = ts
		}
	}
	return myCache, nil
}
