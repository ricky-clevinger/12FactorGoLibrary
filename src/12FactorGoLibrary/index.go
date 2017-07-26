package main

//Author: C Neuhardt
//Last Updated: 7/24/2017

import (
	//"database/sql"
	"html/template"
	"net/http"
	"regexp"
	//"book"
	//"member"
	_ "github.com/go-sql-driver/mysql"
	"member"
)

var validPath = regexp.MustCompile("^/(index.html|admin.html|test.html)$")
var templates = template.Must(template.ParseFiles("views/index.html", "views/admin.html", "views/test.html"))

//Currently not used
type Page struct {
	MemberIds   []int
	MemberFNames []string
}

func loadPage(memberIds []int, memberFNames []string) *Page {
	return &Page{MemberIds: memberIds, MemberFNames: memberFNames}
}

//Renders HTML page
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//Validates path and calls handler
func makeHandler(fn func(http.ResponseWriter, *http.Request, []int, []string), memberIds []int, memberNames []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/*m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}*/
		fn(w, r, memberIds, memberNames)
	}
}

//Handles the index page
func indexHandler(w http.ResponseWriter, r *http.Request, memberIds []int, memberNames []string) {
	p := loadPage(memberIds, memberNames)
	renderTemplate(w, "index", p)
}

//Handles the admin page
func adminHandler(w http.ResponseWriter, r *http.Request, memberIds []int, memberNames []string) {
	p := loadPage(memberIds, memberNames)
	renderTemplate(w, "admin", p)
}

//Handles the test page
func testHandler(w http.ResponseWriter, r *http.Request, memberIds []int, memberNames []string) {
	p := loadPage(memberIds, memberNames)
	renderTemplate(w, "test", p)
}

//Checks for errors
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

//Redirect to index.html
func redirect(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "/index.html", 301)
}

func main() {
	var memberIds []int
	var memberFNames []string

	memberIds = member.GetIds()
	memberFNames = member.GetFNames()

	http.HandleFunc("/", redirect)
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.HandleFunc("/index.html", makeHandler(indexHandler, memberIds, memberFNames))
	http.HandleFunc("/admin.html", makeHandler(adminHandler, memberIds, memberFNames))
	http.HandleFunc("/test.html", makeHandler(testHandler, memberIds, memberFNames))
	http.ListenAndServe(":8080", nil)
}
