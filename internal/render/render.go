package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/shwethadia/HotelReservation/internal/config"
	"github.com/shwethadia/HotelReservation/internal/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

var pathToTemplates = "./templates"

//NewTemplates sets the config for the template package
func NewRenderer(a *config.AppConfig) {

	app = a
}

func AddDefaultData(templateData *models.TemplateData, r *http.Request) *models.TemplateData {

	templateData.Flash = app.Session.PopString(r.Context(), "flash")
	templateData.Error = app.Session.PopString(r.Context(), "error")
	templateData.Warning = app.Session.PopString(r.Context(), "warning")

	templateData.CSRFtoken = nosurf.Token(r)

	return templateData
}

//RenderTemplate renders  templates using html/tsepmlate
func Template(w http.ResponseWriter, r *http.Request, tmp string, templateData *models.TemplateData) error {

	var tc map[string]*template.Template
	if app.UseCache {

		//Get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}
	t, ok := tc[tmp]
	if !ok {

		return errors.New("can't get template from cache")
	}
	buf := new(bytes.Buffer)

	templateData = AddDefaultData(templateData, r)

	_ = t.Execute(buf, templateData)
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
		return err
	}

	return nil

}

func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.htm", pathToTemplates))
	if err != nil {
		return myCache, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.htm", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {

			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.htm", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
