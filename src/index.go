package main

//Author: C Neuhardt
//Last Updated: 7/13/2017

import (
	"html/template"
	"net/http"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
)

var validPath = regexp.MustCompile("^/(index)$")
var templates = template.Must(template.ParseFiles("views/index.html"))

//Currently not used
type Page struct {
	Title string
}

//Renders HTML page
func renderTemplate(w http.ResponseWriter, tmpl string) {
	err := templates.ExecuteTemplate(w, tmpl+".html", nil)
	checkErr(w, err)
}

//Validates path and calls handler
func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r)
	}
}

//Handles the index page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index")
}

//Checks for errors
func checkErr(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.HandleFunc("/index", makeHandler(indexHandler))
	http.ListenAndServe(":8080", nil)
}
