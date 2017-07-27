package main

//Author: C Neuhardt
//Last Updated: 7/26/2017

import (
	"os"
	"database/sql"
	"html/template"
	"net/http"
	"regexp"
	//"book"
	_ "github.com/go-sql-driver/mysql"
	"member"
)

var validPath = regexp.MustCompile("^/(index.html|admin.html|test.html|checkout.html|checkedout)$")
var templates = template.Must(template.ParseFiles("views/index.html", "views/admin.html", "views/test.html", "views/checkout.html"))

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

//Validates path and calls handler
func makeGenericHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/*m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}*/
		fn(w, r)
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

//Handles the checkout page
func checkoutHandler(w http.ResponseWriter, r *http.Request, memberIds []int, memberNames []string) {
	p := loadPage(memberIds, memberNames)
	renderTemplate(w, "checkout", p)
}

//Handles the checkout page
func checkedoutHandler(w http.ResponseWriter, r *http.Request) {
	memberId := r.FormValue("selPerson")
	bookId := r.FormValue("selBook")
	date := r.FormValue("selDateOut")

	db, err := sql.Open("mysql", os.Getenv("LIBRARY"))
	checkErr(err)
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO transaction (book_id, tran_date, che, mid) VALUES (?, ?, 2, ?)")
	checkErr(err)

	stmt.Exec(bookId, date, memberId)

	http.Redirect(w, r, "/index.html", http.StatusFound)
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
	http.HandleFunc("/checkout.html", makeHandler(checkoutHandler, memberIds, memberFNames))
	http.HandleFunc("/checkedout", makeGenericHandler(checkedoutHandler))
	http.ListenAndServe(":8080", nil)
}
