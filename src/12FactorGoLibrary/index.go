package main

//Author: C Neuhardt
//Last Updated: 8/3/2017

import (
	"book"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"helper"
	"html"
	"member"
	"net/http"
	"os"
	"time"
	"handlers"
)


//Handles the index page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	p := handlers.LoadPage(members, books)
	handlers.RenderTemplate(w, "index", p)
}

//Handles the admin page
func adminHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	p := handlers.LoadPage(members, books)
	handlers.RenderTemplate(w, "admin", p)
}

//Handles the books page
func booksHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	members = member.GetMembers()
	books = book.GetBook()

	p := handlers.LoadPage(members, books)
	handlers.RenderTemplate(w, "books", p)
}

//Handles the create book page
func addBookHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	p := handlers.LoadPage(members, books)
	handlers.RenderTemplate(w, "add-book", p)
}

func bookCreatedHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	bookTitle := html.EscapeString(r.FormValue("title"))
	bookAuthF := html.EscapeString(r.FormValue("fName"))
	bookAuthL := html.EscapeString(r.FormValue("lName"))

	if len(bookTitle) == 0 || len(bookAuthF) == 0 || len(bookAuthL) == 0 {
		os.Stderr.WriteString("Empty fields inputted in add-book.html.")
		http.Redirect(w, r, "/add-book.html", http.StatusFound)
	} else {
		book.AddBook(bookTitle, bookAuthF, bookAuthL)

		p := handlers.LoadPage(members, books)
		handlers.RenderTemplate(w, "bookCreated", p)
	}
}

//Handles the edit book page
func editBookHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	id := r.URL.Path[11:]

	books = book.GetBookById(id)

	if len(books) < 1 {
		http.Redirect(w, r, "/books.html", http.StatusFound)
	} else {
		p := handlers.LoadPage(members, books)
		handlers.RenderTemplate(w, "edit-book", p)
	}
}

//Handles the edited book page
func bookEditedHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	bookId := html.EscapeString(r.FormValue("bookId"))
	bookTitle := html.EscapeString(r.FormValue("title"))
	bookAuthF := html.EscapeString(r.FormValue("fName"))
	bookAuthL := html.EscapeString(r.FormValue("lName"))

	if len(bookId) == 0 || len(bookTitle) == 0 || len(bookAuthF) == 0 || len(bookAuthL) == 0 {
		os.Stderr.WriteString("Empty fields inputted in edit-book.html.")
		http.Redirect(w, r, "/edit-book.html", http.StatusFound)
	} else {
		db, err := sql.Open("mysql", os.Getenv("LIBRARY"))
		helper.CheckErr(err)
		defer db.Close()

		stmt, err := db.Prepare("UPDATE books SET book_title = ?, book_authfname = ?, book_authlname = ? WHERE book_id = ?")
		helper.CheckErr(err)

		stmt.Exec(bookTitle, bookAuthF, bookAuthL, bookId)

		p := handlers.LoadPage(members, books)
		handlers.RenderTemplate(w, "bookEdited", p)
	}
}

//Handles the delete book page
func deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	id := r.URL.Path[13:]

	books = book.GetBookById(id)

	if len(books) < 1 {
		http.Redirect(w, r, "/books.html", http.StatusFound)
	} else {
		p := handlers.LoadPage(members, books)
		handlers.RenderTemplate(w, "delete-book", p)
	}
}

//Handles the deleted book page
func bookDeletedHandler(w http.ResponseWriter, r *http.Request) {
	id := html.EscapeString(r.FormValue("bookId"))

	db, err := sql.Open("mysql", os.Getenv("LIBRARY"))
	helper.CheckErr(err)
	defer db.Close()

	//Log transaction
	stmt, err := db.Prepare("DELETE FROM books WHERE book_id = ?")
	helper.CheckErr(err)

	stmt.Exec(id)

	http.Redirect(w, r, "/books.html", http.StatusFound)
}

//Handles the members page
func membersHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	members = member.GetMembers()

	p := handlers.LoadPage(members, books)
	handlers.RenderTemplate(w, "members", p)
}

//Handles the create member page
func addMemberHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	p := handlers.LoadPage(members, books)
	handlers.RenderTemplate(w, "add-member", p)
}

//Handles the create member page
func memberCreatedHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	memberFName := html.EscapeString(r.FormValue("fName"))
	memberLName := html.EscapeString(r.FormValue("lName"))

	if len(memberFName) == 0 || len(memberLName) == 0 {
		os.Stderr.WriteString("Empty fields inputted in add-member.html.")
		http.Redirect(w, r, "/add-member.html", http.StatusFound)
	} else {
		db, err := sql.Open("mysql", os.Getenv("LIBRARY"))
		helper.CheckErr(err)
		defer db.Close()

		stmt, err := db.Prepare("INSERT INTO member (member_fname, member_lname) VALUES (?, ?)")
		helper.CheckErr(err)

		stmt.Exec(memberFName, memberLName)

		p := handlers.LoadPage(members, books)
		handlers.RenderTemplate(w, "memberCreated", p)
	}
}

//Handles the edit member page
func editMemberHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	id := r.URL.Path[13:]

	members = member.GetMemberById(id)

	if len(members) < 1 {
		http.Redirect(w, r, "/members.html", http.StatusFound)
	} else {
		p := handlers.LoadPage(members, books)
		handlers.RenderTemplate(w, "edit-member", p)
	}
}

//Handles the edited member page
func memberEditedHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	memberId := html.EscapeString(r.FormValue("memId"))
	memberFName := html.EscapeString(r.FormValue("fName"))
	memberLName := html.EscapeString(r.FormValue("lName"))

	if len(memberId) == 0 || len(memberFName) == 0 || len(memberLName) == 0 {
		os.Stderr.WriteString("Empty fields inputted in edit-member.html.")
		http.Redirect(w, r, "/edit-member.html", http.StatusFound)
	} else {
		db, err := sql.Open("mysql", os.Getenv("LIBRARY"))
		helper.CheckErr(err)
		defer db.Close()

		//Log transaction
		stmt, err := db.Prepare("UPDATE member SET member_fname = ?, member_lname = ? WHERE member_id = ?")
		helper.CheckErr(err)

		stmt.Exec(memberFName, memberLName, memberId)

		p := handlers.LoadPage(members, books)
		handlers.RenderTemplate(w, "memberEdited", p)
	}
}

//Handles the delete member page
func deleteMemberHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	id := r.URL.Path[15:]

	members = member.GetMemberById(id)

	if len(members) < 1 {
		http.Redirect(w, r, "/members.html", http.StatusFound)
	} else {
		p := handlers.LoadPage(members, books)
		handlers.RenderTemplate(w, "delete-member", p)
	}
}

//Hanldes the deleted member page
func memberDeletedHandler(w http.ResponseWriter, r *http.Request) {
	id := html.EscapeString(r.FormValue("memId"))

	db, err := sql.Open("mysql", os.Getenv("LIBRARY"))
	helper.CheckErr(err)
	defer db.Close()

	//Log transaction
	stmt, err := db.Prepare("DELETE FROM member WHERE member_id = ?")
	helper.CheckErr(err)

	stmt.Exec(id)

	http.Redirect(w, r, "/members.html", http.StatusFound)
}

//Handles the test page
func testHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	members = member.GetMembers()
	books = book.GetBook()

	p := handlers.LoadPage(members, books)
	handlers.RenderTemplate(w, "test", p)
}

//Handles the checkout page
func checkoutHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	members = member.GetMembers()

	var books []book.Book
	books = book.GetCheckedInBook()

	p := handlers.LoadPage(members, books)
	handlers.RenderTemplate(w, "checkout", p)
}

//Handles the checkout page
func checkedoutHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	current_time := time.Now().Local()

	memberId := html.EscapeString(r.FormValue("selPerson"))
	bookId := html.EscapeString(r.FormValue("selBook"))
	date := current_time.Format("2006-01-02 15:04:05")

	db, err := sql.Open("mysql", os.Getenv("LIBRARY"))
	helper.CheckErr(err)
	defer db.Close()

	//Log transaction
	stmt, err := db.Prepare("INSERT INTO transaction (book_id, tran_date, che, mid) VALUES (?, ?, 2, ?)")
	helper.CheckErr(err)

	stmt.Exec(bookId, date, memberId)

	//Update checkout status
	stmt2, err := db.Prepare("UPDATE books SET book_check=2, mid=?, book_out_date=? WHERE book_id=? AND book_check = 1")
	helper.CheckErr(err)

	stmt2.Exec(memberId, date, bookId)

	p := handlers.LoadPage(members, books)
	handlers.RenderTemplate(w, "checkedout", p)
}

//Handles the checkin page
func checkinHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	members = member.GetMembers()

	var books []book.Book
	books = book.GetCheckedOutBook()

	p := handlers.LoadPage(members, books)
	handlers.RenderTemplate(w, "checkin", p)
}

//Handles the checkin page
func checkedinHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	current_time := time.Now().Local()

	bookId := html.EscapeString(r.FormValue("selBook"))
	date := current_time.Format("2006-01-02 15:04:05")

	db, err := sql.Open("mysql", os.Getenv("LIBRARY"))
	helper.CheckErr(err)
	defer db.Close()

	//Log transaction
	stmt, err := db.Prepare("INSERT INTO transaction (book_id, tran_date, che, mid) VALUES (?, ?, 1, (SELECT mid FROM books WHERE book_id = ?))")
	helper.CheckErr(err)

	stmt.Exec(bookId, date, bookId)

	//Update checkout status
	stmt2, err := db.Prepare("UPDATE books SET book_check=1, mid=0, book_out_date=null WHERE book_id=? AND book_check = 2")
	helper.CheckErr(err)

	stmt2.Exec(bookId)

	p := handlers.LoadPage(members, books)
	handlers.RenderTemplate(w, "checkedin", p)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {

	var books []book.Book
	var members []member.Member
	search := html.EscapeString(r.FormValue("s-bar"))

	if len(search) < 1 {
		os.Stderr.WriteString("Empty fields inputted in home page search.")
		http.Redirect(w, r, "/index.html", http.StatusFound)
	} else {
		books = book.GetSearchedBook(search)
		members = member.GetSearchedMember(search)

		p := handlers.LoadPage(members, books)
		handlers.RenderTemplate(w, "results", p)
	}
}

/*

//Checks for errors
func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

*/

//Redirect to index.html
func redirect(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "/index.html", 301)
}

//Handles
func handles() {

	http.HandleFunc("/index.html", handlers.MakeHandler(indexHandler))
	http.HandleFunc("/search", handlers.MakeHandler(searchHandler))
	http.HandleFunc("/admin.html", handlers.MakeHandler(adminHandler))
	http.HandleFunc("/books.html", handlers.MakeHandler(booksHandler))
	http.HandleFunc("/add-book.html", handlers.MakeHandler(addBookHandler))
	http.HandleFunc("/bookCreated", handlers.MakeHandler(bookCreatedHandler))
	http.HandleFunc("/edit-book/", handlers.MakeHandler(editBookHandler))
	http.HandleFunc("/bookEdited", handlers.MakeHandler(bookEditedHandler))
	http.HandleFunc("/delete-book/", handlers.MakeHandler(deleteBookHandler))
	http.HandleFunc("/bookDeleted", handlers.MakeHandler(bookDeletedHandler))
	http.HandleFunc("/members.html", handlers.MakeHandler(membersHandler))
	http.HandleFunc("/add-member.html", handlers.MakeHandler(addMemberHandler))
	http.HandleFunc("/memberCreated", handlers.MakeHandler(memberCreatedHandler))
	http.HandleFunc("/edit-member/", handlers.MakeHandler(editMemberHandler))
	http.HandleFunc("/memberEdited", handlers.MakeHandler(memberEditedHandler))
	http.HandleFunc("/delete-member/", handlers.MakeHandler(deleteMemberHandler))
	http.HandleFunc("/memberDeleted", handlers.MakeHandler(memberDeletedHandler))
	http.HandleFunc("/test.html", handlers.MakeHandler(testHandler))
	http.HandleFunc("/checkout.html", handlers.MakeHandler(checkoutHandler))
	http.HandleFunc("/checkedout", handlers.MakeHandler(checkedoutHandler))
	http.HandleFunc("/checkin.html", handlers.MakeHandler(checkinHandler))
	http.HandleFunc("/checkedin", handlers.MakeHandler(checkedinHandler))
}

func main() {
	http.HandleFunc("/", redirect)
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	handles()
	http.ListenAndServe(":8080", nil)
}
