package main

//Author: C Neuhardt
//Last Updated: 7/17/2017

import (
	"database/sql"
	"html/template"
	"net/http"
	"regexp"

	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var validPath = regexp.MustCompile("^/(index)$")
var templates = template.Must(template.ParseFiles("views/index.html"))

/*Currently not used
type Page struct {
	member *Rows
}*/

//Renders HTML page
func renderTemplate(w http.ResponseWriter, tmpl string) {
	err := templates.ExecuteTemplate(w, tmpl+".html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func main() {
	//TODO: Replace connection data with env vars
	db, err := sql.Open("mysql", "cgidevmem:Password1@tcp(cgiprojdevmember.cxyeb3wmov3g.us-east-1.rds.amazonaws.com:3306)/cgiprojdevmember")
	checkErr(err)
	defer db.Close()
	fmt.Println("db connected.")

	memberRows, err := db.Query("SELECT * FROM member")
	checkErr(err)
	fmt.Println("Rows Acquired.")

	for memberRows.Next() {
		var MemberId int
		var MemberFName string
		var MemberLName string
		err = memberRows.Scan(&MemberId, &MemberFName, &MemberLName)
		checkErr(err)
		fmt.Println(MemberId)
		fmt.Println(MemberFName)
		fmt.Println(MemberLName)
	}

	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.HandleFunc("/index", makeHandler(indexHandler))
	http.ListenAndServe(":8080", nil)
}
