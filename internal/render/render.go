package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/justinas/nosurf"
	"github.com/shwethadia/HotelReservation/internal/config"
	"github.com/shwethadia/HotelReservation/internal/models"
)

var functions = template.FuncMap{

	"humanDate":  HumanDate,
	"formatDate": FormatDate,
	"iterate":    Iterate,
	"add":        Add,
}

var app *config.AppConfig

var pathToTemplates = "./templates"

func Add(a, b int) int {

	return a + b
}

//Iterate returns a slice of ints, starting at 1, going to count
func Iterate(count int) []int {

	var i int
	var items []int
	for i = 0; i < count; i++ {

		items = append(items, i)
	}
	return items
}

//NewTemplates sets the config for the template package
func NewRenderer(a *config.AppConfig) {

	app = a
}

//HumanDate Returns time in yyyy-mm-dd format
func HumanDate(t time.Time) string {

	return t.Format("2006-01-02")
}

func FormatDate(t time.Time, f string) string {

	return t.Format(f)

}

func AddDefaultData(templateData *models.TemplateData, r *http.Request) *models.TemplateData {

	templateData.Flash = app.Session.PopString(r.Context(), "flash")
	templateData.Error = app.Session.PopString(r.Context(), "error")
	templateData.Warning = app.Session.PopString(r.Context(), "warning")

	templateData.CSRFtoken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {

		templateData.IsAuthenticated = 1
	}

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
