package main

//Author: C Neuhardt
//Last Updated: 8/3/2017

import (
	"os"
	"fmt"
	"database/sql"
	"html/template"
	"net/http"
	"regexp"
	"book"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"member"
	"time"
)

var validPath = regexp.MustCompile("^/(index.html|search|results.html|admin.html|books.html|add-book.html|bookCreated|edit-book/[0-9]+|bookEdited|delete-book/[0-9]+|bookDeleted|members.html|add-member.html|memberCreated|edit-member/[0-9]+|memberEdited|delete-member/[0-9]+|memberDeleted|test.html|checkout.html|checkedout|checkin.html|checkedin)$")
var templates = template.Must(template.ParseFiles("views/index.html", "views/admin.html", "views/books.html", "views/add-book.html", "views/bookCreated.html", "views/edit-book.html", "views/bookEdited.html", "views/delete-book.html", "views/members.html", "views/add-member.html", "views/memberCreated.html", "views/edit-member.html", "views/memberEdited.html", "views/delete-member.html", "views/test.html", "views/checkout.html", "views/checkedout.html", "views/checkin.html", "views/checkedin.html", "views/results.html"))

type Page struct {
	Members []member.Member
	Books []book.Book
}

func loadPage(members []member.Member, books []book.Book) *Page {
	if(len(members) > 0) {
		fmt.Println("Loaded member #: ", len(members))
	}
	if(len(books) > 0) {
		fmt.Println("Loaded book #: ", len(books))
	}
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

//Handles the create book page
func addBookHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book
	
	p := loadPage(members, books)
	renderTemplate(w, "add-book", p)
}

func bookCreatedHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	bookTitle := r.FormValue("title")
	bookAuthF := r.FormValue("fName")
	bookAuthL := r.FormValue("lName")

	if len(bookTitle)==0 || len(bookAuthF)==0 || len(bookAuthL)==0 {
		os.Stderr.WriteString("Empty fields inputted in add-book.html.")
		http.Redirect(w, r, "/add-book.html", http.StatusFound)
	} else {
		book.AddBook(bookTitle, bookAuthF, bookAuthL)

		p := loadPage(members, books)
		renderTemplate(w, "bookCreated", p)
	}
}

//Handles the edit book page
func editBookHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	id := r.URL.Path[11:]

	books = book.GetBookById(id)
	
	if(len(books) < 1) {
		http.Redirect(w, r, "/books.html", http.StatusFound)
	} else {
		p := loadPage(members, books)
		renderTemplate(w, "edit-book", p)
	}
}

//Handles the edited book page
func bookEditedHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book
	
	bookId := r.FormValue("bookId")
	bookTitle := r.FormValue("title")
	bookAuthF := r.FormValue("fName")
	bookAuthL := r.FormValue("lName")

	if len(bookId)==0 || len(bookTitle)==0 || len(bookAuthF)==0 || len(bookAuthL)==0 {
		os.Stderr.WriteString("Empty fields inputted in edit-book.html.")
		http.Redirect(w, r, "/edit-book.html", http.StatusFound)
	} else {
		db, err := sql.Open("mysql", os.Getenv("LIBRARY"))
		checkErr(err)
		defer db.Close()

		stmt, err := db.Prepare("UPDATE books SET book_title = ?, book_authfname = ?, book_authlname = ? WHERE book_id = ?")
		checkErr(err)

		stmt.Exec(bookTitle, bookAuthF, bookAuthL, bookId)

		p := loadPage(members, books)
		renderTemplate(w, "bookEdited", p)
	}
}

//Handles the delete book page
func deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book
	
	id := r.URL.Path[13:]
	
	books = book.GetBookById(id)

	if(len(books) < 1) {
		http.Redirect(w, r, "/books.html", http.StatusFound)
	} else {
		p := loadPage(members, books)
		renderTemplate(w, "delete-book", p)
	}
}
//Handles the deleted book page
func bookDeletedHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("bookId")
	
	db, err := sql.Open("mysql", os.Getenv("LIBRARY"))
	checkErr(err)
	defer db.Close()

	//Log transaction
	stmt, err := db.Prepare("DELETE FROM books WHERE book_id = ?")
	checkErr(err)

	stmt.Exec(id)

	http.Redirect(w, r, "/books.html", http.StatusFound)
}

//Handles the members page
func membersHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book
	
	members = member.GetMembers()

	p := loadPage(members, books)
	renderTemplate(w, "members", p)
}

//Handles the create member page
func addMemberHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book
	
	p := loadPage(members, books)
	renderTemplate(w, "add-member", p)
}

//Handles the create member page
func memberCreatedHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	memberFName := r.FormValue("fName")
	memberLName := r.FormValue("lName")
	
	if len(memberFName)==0 || len(memberLName)==0 {
		os.Stderr.WriteString("Empty fields inputted in add-member.html.")
		http.Redirect(w, r, "/add-member.html", http.StatusFound)
	} else {
		db, err := sql.Open("mysql", os.Getenv("LIBRARY"))
		checkErr(err)
		defer db.Close()

		stmt, err := db.Prepare("INSERT INTO member (member_fname, member_lname) VALUES (?, ?)")
		checkErr(err)

		stmt.Exec(memberFName, memberLName)

		p := loadPage(members, books)
		renderTemplate(w, "memberCreated", p)
	}
}

//Handles the edit member page
func editMemberHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	id := r.URL.Path[13:]

	members = member.GetMemberById(id)
	
	if(len(members) < 1) {
		http.Redirect(w, r, "/members.html", http.StatusFound)
	} else {
		p := loadPage(members, books)
		renderTemplate(w, "edit-member", p)
	}
}

//Handles the edited member page
func memberEditedHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book
	
	memberId := r.FormValue("memId")
	memberFName := r.FormValue("fName")
	memberLName := r.FormValue("lName")
	
	if len(memberId)==0 || len(memberFName)==0 || len(memberLName)==0 {
		os.Stderr.WriteString("Empty fields inputted in edit-member.html.")
		http.Redirect(w, r, "/edit-member.html", http.StatusFound)
	} else {
		db, err := sql.Open("mysql", os.Getenv("LIBRARY"))
		checkErr(err)
		defer db.Close()

		//Log transaction
		stmt, err := db.Prepare("UPDATE member SET member_fname = ?, member_lname = ? WHERE member_id = ?")
		checkErr(err)

		stmt.Exec(memberFName, memberLName, memberId)

		p := loadPage(members, books)
		renderTemplate(w, "memberEdited", p)
	}
}

//Handles the delete member page
func deleteMemberHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book
	
	id := r.URL.Path[15:]
	
	members = member.GetMemberById(id)

	if(len(members) < 1) {
		http.Redirect(w, r, "/members.html", http.StatusFound)
	} else {
		p := loadPage(members, books)
		renderTemplate(w, "delete-member", p)
	}
}

//Hanldes the deleted member page
func memberDeletedHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("memId")
	
	db, err := sql.Open("mysql", os.Getenv("LIBRARY"))
	checkErr(err)
	defer db.Close()

	//Log transaction
	stmt, err := db.Prepare("DELETE FROM member WHERE member_id = ?")
	checkErr(err)

	stmt.Exec(id)

	http.Redirect(w, r, "/members.html", http.StatusFound)
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
	var members []member.Member
	var books []book.Book
	
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

	p := loadPage(members, books)
	renderTemplate(w, "checkedout", p)
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
	var members []member.Member
	var books []book.Book
	
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

	p := loadPage(members, books)
	renderTemplate(w, "checkedin", p)
}

func searchHandler(w http.ResponseWriter, r *http.Request){
	
	var books []book.Book
	var members []member.Member
	search := r.FormValue("s-bar")
	books = book.GetSearchedBook(search)
	members = member.GetSearchedMember(search)

	p := loadPage(members, books)
	renderTemplate(w, "results", p)
	
}

//Checks for errors
func checkErr(err error) {
	if err != nil {
		log.Panic(err)
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
	http.HandleFunc("/add-book.html", makeHandler(addBookHandler))
	http.HandleFunc("/bookCreated", makeHandler(bookCreatedHandler))
	http.HandleFunc("/edit-book/", makeHandler(editBookHandler))
	http.HandleFunc("/bookEdited", makeHandler(bookEditedHandler))
	http.HandleFunc("/delete-book/", makeHandler(deleteBookHandler))
	http.HandleFunc("/bookDeleted", makeHandler(bookDeletedHandler))
	http.HandleFunc("/members.html", makeHandler(membersHandler))
	http.HandleFunc("/add-member.html", makeHandler(addMemberHandler))
	http.HandleFunc("/memberCreated", makeHandler(memberCreatedHandler))
	http.HandleFunc("/edit-member/", makeHandler(editMemberHandler))
	http.HandleFunc("/memberEdited", makeHandler(memberEditedHandler))
	http.HandleFunc("/delete-member/", makeHandler(deleteMemberHandler))
	http.HandleFunc("/memberDeleted", makeHandler(memberDeletedHandler))
	http.HandleFunc("/test.html", makeHandler(testHandler))
	http.HandleFunc("/checkout.html", makeHandler(checkoutHandler))
	http.HandleFunc("/checkedout", makeHandler(checkedoutHandler))
	http.HandleFunc("/checkin.html", makeHandler(checkinHandler))
	http.HandleFunc("/checkedin", makeHandler(checkedinHandler))
	http.ListenAndServe(":8080", nil)
}
