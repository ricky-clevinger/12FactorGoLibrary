package main

//Author: C Neuhardt
//Last Updated: 7/17/2017

import (
	"database/sql"
	"html/template"
	"net/http"
	"regexp"
	_"vendor/github.com/go-sql-driver/mysql"
)

var validPath = regexp.MustCompile("^/(index.html|admin.html|test.html)$")
var templates = template.Must(template.ParseFiles("views/index.html", "views/admin.html", "views/test.html"))

//Currently not used
type Page struct {
	MemberIds   []int
	MemberNames []string
}

func loadPage(memberIds []int, memberNames []string) *Page {
	return &Page{MemberIds: memberIds, MemberNames: memberNames}
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

func main() {
	var memberIds []int
	var memberNames []string

	//TODO: Replace connection data with env vars
	db, err := sql.Open("mysql", "cgidevmem:Password1@tcp(cgiprojdevmember.cxyeb3wmov3g.us-east-1.rds.amazonaws.com:3306)/cgiprojdevmember")
	checkErr(err)
	defer db.Close()

	memberIdRows, err := db.Query("SELECT member_id FROM member")
	checkErr(err)
	for memberIdRows.Next() {
		var MemberId int
		err = memberIdRows.Scan(&MemberId)
		checkErr(err)
		memberIds = append(memberIds, MemberId)
	}

	memberNameRows, err := db.Query("SELECT member_fname, member_lname FROM member")
	checkErr(err)
	for memberNameRows.Next() {
		var MemberFName string
		var MemberLName string
		err = memberNameRows.Scan(&MemberFName, &MemberLName)
		checkErr(err)
		memberNames = append(memberNames, MemberFName+" "+MemberLName)
	}

	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.HandleFunc("/index.html", makeHandler(indexHandler, memberIds, memberNames))
	http.HandleFunc("/admin.html", makeHandler(adminHandler, memberIds, memberNames))
	http.HandleFunc("/test.html", makeHandler(testHandler, memberIds, memberNames))
	http.ListenAndServe(":8080", nil)
}
