package main

import (
	"database/sql"
	"fmt"
	"github/chiragg6/url-shortner/internals/models"
	"log"
	"net/http"
	"text/template"

	"golang.org/x/text/internal/number"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var functions = template.FuncMap{
	"formatClicks": formatClicks,
}

type PageData struct {
	BaseURL, Error string
	URLData        []*models.ShortnerData
}

type App struct {
	urls *models.ShortnerDataModel
}

func serverError(w http.ResponseWriter, err error) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func newApp(dbFile string) App {
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	return App{urls: &models.ShortnerDataModel{DB: db}}
}

func formatClicks(clicks int) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%v", number.Decimal(clicks))
}

func (a *App) getDefaultRoute(w http.ResponseWriter, r *http.Request) {
	tmpFile := "./templates/default.html"
	tmpl, err := template.New("default.html").Funcs(functions).ParseFiles(tmpFile)
	if err != nil {
		fmt.Println(err.Error())
		serverError(w, err)
		return
	}

	urls, err := a.urls.Latest()
	if err != nil {
		fmt.Printf("Could not retrieve all URLs, because %s.\n", err)
		return
	}

	baseURL := "http://" + r.Host + "/"
	pageData := PageData{
		URLData: urls,
		BaseURL: baseURL,
	}

	err = tmpl.Execute(w, pageData)
	if err != nil {
		fmt.Println(err.Error())
		serverError(w, err)
	}
}

func (a *App) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./static/"))
}
