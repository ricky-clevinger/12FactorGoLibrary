package main

//Author: C Neuhardt
//Last Updated: 8/1/2017

import (
	"os"
	//"fmt"
	"database/sql"
	"html/template"
	"net/http"
	"regexp"
	"book"

	_ "github.com/go-sql-driver/mysql"
	"member"
	"time"
)

var validPath = regexp.MustCompile("^/(index.html|search|results.html|admin.html|books.html|members.html|test.html|checkout.html|checkedout|checkin.html|checkedin)$")
var templates = template.Must(template.ParseFiles("views/index.html", "views/admin.html", "views/books.html", "views/members.html", "views/test.html", "views/checkout.html", "views/checkin.html", "views/results.html"))

type Page struct {
	Members []member.Member
	Books []book.Book
}

func loadPage(members []member.Member, books []book.Book) *Page {
	return &Page{Members:members, Books:books}
}

//Renders HTML page
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//Validates path and calls handler
func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
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
func indexHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book
	
	p := loadPage(members, books)
	renderTemplate(w, "index", p)
}

//Handles the admin page
func adminHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book
	
	p := loadPage(members, books)
	renderTemplate(w, "admin", p)
}

//Handles the books page
func booksHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book
	
	members = member.GetMembers()
	books = book.GetBook()

	p := loadPage(members, books)
	renderTemplate(w, "books", p)
}

//Handles the members page
func membersHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book
	
	members = member.GetMembers()

	p := loadPage(members, books)
	renderTemplate(w, "members", p)
}

//Handles the test page
func testHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book
	
	members = member.GetMembers()
	books = book.GetBook()
	
	p := loadPage(members, books)
	renderTemplate(w, "test", p)
}

//Handles the checkout page
func checkoutHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	members = member.GetMembers()
	
	var books []book.Book
	books = book.GetCheckedInBook()
	
	p := loadPage(members, books)
	renderTemplate(w, "checkout", p)
}

//Handles the checkout page
func checkedoutHandler(w http.ResponseWriter, r *http.Request) {
	current_time := time.Now().Local()
	
	memberId := r.FormValue("selPerson")
	bookId := r.FormValue("selBook")
	date := current_time.Format("2006-01-02 15:04:05")
	
	db, err := sql.Open("mysql", os.Getenv("LIBRARY"))
	checkErr(err)
	defer db.Close()

	//Log transaction
	stmt, err := db.Prepare("INSERT INTO transaction (book_id, tran_date, che, mid) VALUES (?, ?, 2, ?)")
	checkErr(err)

	stmt.Exec(bookId, date, memberId)

	//Update checkout status
	stmt2, err := db.Prepare("UPDATE books SET book_check=2, mid=?, book_out_date=? WHERE book_id=? AND book_check = 1")
	checkErr(err)

	stmt2.Exec(memberId, date, bookId)

	http.Redirect(w, r, "/index.html", http.StatusFound)
}

//Handles the checkin page
func checkinHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	members = member.GetMembers()
	
	var books []book.Book
	books = book.GetCheckedOutBook()
	
	p := loadPage(members, books)
	renderTemplate(w, "checkin", p)
}

//Handles the checkin page
func checkedinHandler(w http.ResponseWriter, r *http.Request) {
	current_time := time.Now().Local()
	
	bookId := r.FormValue("selBook")
	date := current_time.Format("2006-01-02 15:04:05")

	db, err := sql.Open("mysql", os.Getenv("LIBRARY"))
	checkErr(err)
	defer db.Close()

	//Log transaction
	stmt, err := db.Prepare("INSERT INTO transaction (book_id, tran_date, che, mid) VALUES (?, ?, 1, (SELECT mid FROM books WHERE book_id = ?))")
	checkErr(err)

	stmt.Exec(bookId, date, bookId)

	//Update checkout status
	stmt2, err := db.Prepare("UPDATE books SET book_check=1, mid=0, book_out_date=null WHERE book_id=? AND book_check = 2")
	checkErr(err)

	stmt2.Exec(bookId)

	http.Redirect(w, r, "/index.html", http.StatusFound)
}

func searchHandler(w http.ResponseWriter, r *http.Request){
	
	var books []book.Book
	var members []member.Member
	search := r.FormValue("s-bar")
	books = book.GetSearchedBook(search)

	p := loadPage(members, books)
	renderTemplate(w, "results", p)
	
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
	http.HandleFunc("/", redirect)
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.HandleFunc("/index.html", makeHandler(indexHandler))
	http.HandleFunc("/search", makeHandler(searchHandler))
	http.HandleFunc("/admin.html", makeHandler(adminHandler))
	http.HandleFunc("/books.html", makeHandler(booksHandler))
	http.HandleFunc("/members.html", makeHandler(membersHandler))
	http.HandleFunc("/test.html", makeHandler(testHandler))
	http.HandleFunc("/checkout.html", makeHandler(checkoutHandler))
	http.HandleFunc("/checkedout", makeHandler(checkedoutHandler))
	http.HandleFunc("/checkin.html", makeHandler(checkinHandler))
	http.HandleFunc("/checkedin", makeHandler(checkedinHandler))
	http.ListenAndServe(":8080", nil)
}
