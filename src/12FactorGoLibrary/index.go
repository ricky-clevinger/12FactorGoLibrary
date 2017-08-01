package main

//Author: C Neuhardt
//Last Updated: 8/1/2017

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
/*type Page struct {
	Member_id []int
	Member_fname []string
	Member_lname []string
	Book_id []int
	Book_title []string
	Book_authF []string
	Book_authL []string
	Library_id []int
	Book_check []int
	Mid []int
	Book_out_date []string
}*/

type Page struct {
	Members []member.Member
	Books []book.Book
}

func loadPage(members []member.Member, books []book.Book) *Page {
	//return &Page{Member_id:member_ids, Member_fname:member_fnames, Member_lname:member_lnames, Book_id:book_ids, Book_title:book_titles, Book_authF:book_authFs, Book_authL:book_authLs, Library_id:lib_ids, Book_check:book_checks, Mid:mids, Book_out_date:book_out_dates}
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
func makeHandler(fn func(http.ResponseWriter, *http.Request, []member.Member, []book.Book), members []member.Member, books []book.Book) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/*m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}*/
		fn(w, r, members, books)
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
func indexHandler(w http.ResponseWriter, r *http.Request, members []member.Member, books []book.Book) {
	p := loadPage(members, books)
	renderTemplate(w, "index", p)
}

//Handles the admin page
func adminHandler(w http.ResponseWriter, r *http.Request, members []member.Member, books []book.Book) {
	p := loadPage(members, books)
	renderTemplate(w, "admin", p)
}

//Handles the test page
func testHandler(w http.ResponseWriter, r *http.Request, members []member.Member, books []book.Book) {
	p := loadPage(members, books)
	renderTemplate(w, "test", p)
}

//Handles the checkout page
func checkoutHandler(w http.ResponseWriter, r *http.Request, members []member.Member, books []book.Book) {
	p := loadPage(members, books)
	renderTemplate(w, "checkout", p)
}

//Handles the checkout page
func checkedoutHandler(w http.ResponseWriter, r *http.Request) {
	memberId := r.FormValue("selPerson")
	bookId := r.FormValue("selBook")
	date := r.FormValue("selDateOut")
	
	fmt.Println(memberId, bookId, date)
	
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
func checkinHandler(w http.ResponseWriter, r *http.Request, members []member.Member, books []book.Book) {
	p := loadPage(members, books)
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
	var members []member.Member
	var books []book.Book

	members = member.GetMembers()
	books = book.GetBook()

	http.HandleFunc("/", redirect)
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.HandleFunc("/index.html", makeHandler(indexHandler, members, books))
	http.HandleFunc("/admin.html", makeHandler(adminHandler, members, books))
	http.HandleFunc("/test.html", makeHandler(testHandler, members, books))
	http.HandleFunc("/checkout.html", makeHandler(checkoutHandler, members, books))
	http.HandleFunc("/checkedout", makeGenericHandler(checkedoutHandler))
	http.HandleFunc("/checkin.html", makeHandler(checkinHandler, members, books))
	http.HandleFunc("/checkedin", makeGenericHandler(checkedinHandler))
	http.ListenAndServe(":8080", nil)
}
