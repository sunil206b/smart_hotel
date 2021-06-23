package render

import (
	"bytes"
	"errors"
	"github.com/justinas/nosurf"
	"github.com/sunil206b/smart_booking/internal/config"
	"github.com/sunil206b/smart_booking/internal/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}

var appConfig *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	appConfig = a
}

// RenderTemplate function will load the html files from the specified location and parses
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data *models.TemplateData) {
	var tc map[string]*template.Template
	if appConfig.UseCache {
		tc = appConfig.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}
	t, ok := tc[tmpl]
	if !ok {
		log.Fatalf("Template %v not exisat in the template cache\n", tmpl)
	}
	buf := new(bytes.Buffer)
	AddDefaultData(data, r)
	err := t.Execute(buf, data)
	if err != nil {
		log.Printf("Failed to read data from cache buffer:  %v\n", err.Error())
		return
	}
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Printf("Failed to write template to the browser:  %v\n", err.Error())
	}
}

// CreateTemplateCache function creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return nil, errors.New("Error while looking for pages " + err.Error())
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, errors.New("Error while generating template set " + err.Error())
		}
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return nil, errors.New("Error while checking for layout file " + err.Error())
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return nil, errors.New("Error while parsing layout file " + err.Error())
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}

func AddDefaultData(data *models.TemplateData, r *http.Request) {
	data.CSRFToken = nosurf.Token(r)
}