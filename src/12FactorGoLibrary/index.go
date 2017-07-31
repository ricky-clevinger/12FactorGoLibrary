package main

//Author: C Neuhardt
//Last Updated: 7/26/2017

import (
	"os"
	"fmt"
	"database/sql"
	"html/template"
	"net/http"
	"regexp"
	"book"

	_ "github.com/go-sql-driver/mysql"
	"member"
)

var validPath = regexp.MustCompile("^/(index.html|admin.html|test.html|checkout.html|checkedout|checkin.html|checkedin)$")
var templates = template.Must(template.ParseFiles("views/index.html", "views/admin.html", "views/test.html", "views/checkout.html", "views/checkin.html"))

//Currently not used
type Page struct {
	MemberIds   []int
	MemberFNames []string
	Books []book.Book
}

func loadPage(memberIds []int, memberFNames []string, books []book.Book) *Page {
	return &Page{MemberIds: memberIds, MemberFNames: memberFNames, Books:books}
}

//Renders HTML page
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//Validates path and calls handler
func makeHandler(fn func(http.ResponseWriter, *http.Request, []int, []string, []book.Book), memberIds []int, memberNames []string, books []book.Book) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/*m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}*/
		fn(w, r, memberIds, memberNames, books)
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
func indexHandler(w http.ResponseWriter, r *http.Request, memberIds []int, memberNames []string, books []book.Book) {
	p := loadPage(memberIds, memberNames, books)
	renderTemplate(w, "index", p)
}

//Handles the admin page
func adminHandler(w http.ResponseWriter, r *http.Request, memberIds []int, memberNames []string, books []book.Book) {
	p := loadPage(memberIds, memberNames, books)
	renderTemplate(w, "admin", p)
}

//Handles the test page
func testHandler(w http.ResponseWriter, r *http.Request, memberIds []int, memberNames []string, books []book.Book) {
	p := loadPage(memberIds, memberNames, books)
	renderTemplate(w, "test", p)
}

//Handles the checkout page
func checkoutHandler(w http.ResponseWriter, r *http.Request, memberIds []int, memberNames []string, books []book.Book) {
	p := loadPage(memberIds, memberNames, books)
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

	//Log transaction
	stmt, err := db.Prepare("INSERT INTO transaction (book_id, tran_date, che, mid) VALUES (?, ?, 2, ?)")
	checkErr(err)

	stmt.Exec(bookId, date, memberId)

	//Update checkout status
	stmt2, err := db.Prepare("UPDATE books SET book_check=2, mid=?, book_out_date=? WHERE book_id=?")
	checkErr(err)

	stmt2.Exec(memberId, date, bookId)

	http.Redirect(w, r, "/index.html", http.StatusFound)
}

//Handles the checkin page
func checkinHandler(w http.ResponseWriter, r *http.Request, memberIds []int, memberNames []string, books []book.Book) {
	p := loadPage(memberIds, memberNames, books)
	renderTemplate(w, "checkin", p)
}

//Handles the checkin page
func checkedinHandler(w http.ResponseWriter, r *http.Request) {
	var memberId int
	bookId := r.FormValue("selBook")
	date := r.FormValue("selDateIn")

	db, err := sql.Open("mysql", os.Getenv("LIBRARY"))
	checkErr(err)
	defer db.Close()

	//Get Member ID
	queryString := fmt.Sprintf("SELECT mid FROM books WHERE book_id = %d", bookId)
	rows, err := db.Query(queryString)
	checkErr(err)
	
	for rows.Next() {
		var mid int
		err = rows.Scan(&mid)
		checkErr(err)
		memberId = mid
	}

	//Log transaction
	stmt, err := db.Prepare("INSERT INTO transaction (book_id, tran_date, che, mid) VALUES (?, ?, 2, ?)")
	checkErr(err)

	stmt.Exec(bookId, date, memberId)

	//Update checkout status
	stmt2, err := db.Prepare("UPDATE books SET book_check=1, mid=0, book_out_date=null WHERE book_id=?")
	checkErr(err)

	stmt2.Exec(bookId)

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
	var books []book.Book

	memberIds = member.GetIds()
	memberFNames = member.GetFNames()
	books = book.GetBook()

	//fmt.Println(books)

	http.HandleFunc("/", redirect)
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.HandleFunc("/index.html", makeHandler(indexHandler, memberIds, memberFNames, books))
	http.HandleFunc("/admin.html", makeHandler(adminHandler, memberIds, memberFNames, books))
	http.HandleFunc("/test.html", makeHandler(testHandler, memberIds, memberFNames, books))
	http.HandleFunc("/checkout.html", makeHandler(checkoutHandler, memberIds, memberFNames, books))
	http.HandleFunc("/checkedout", makeGenericHandler(checkedoutHandler))
	http.HandleFunc("/checkin.html", makeHandler(checkinHandler, memberIds, memberFNames, books))
	http.HandleFunc("/checkedin", makeGenericHandler(checkedinHandler))
	http.ListenAndServe(":8080", nil)
}
