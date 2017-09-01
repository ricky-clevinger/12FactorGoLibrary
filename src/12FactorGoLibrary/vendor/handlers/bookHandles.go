package handlers

//Created By: Ricky Clevinger
//Updated On: 8/17/2017
//Last Updated By: Ricky Clevinger

import (
	"book"
	"database/sql"
	_ "webPackages/github.com/go-sql-driver/mysql"
	"helper"
	"member"
	"net/http"
	"os"
)

//Handles the books page
func BooksHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	members = member.GetMembers()
	books = book.GetBook()

	p := LoadPage(members, books)
	RenderTemplate(w, "books", p)
}

//Handles the create book page
func AddBookHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	p := LoadPage(members, books)
	RenderTemplate(w, "add-book", p)
}

func BookCreatedHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	bookTitle := helper.HTMLClean(r.FormValue("title"))
	bookAuthF := helper.HTMLClean(r.FormValue("fName"))
	bookAuthL := helper.HTMLClean(r.FormValue("lName"))

	if len(bookTitle) == 0 || len(bookAuthF) == 0 || len(bookAuthL) == 0 {
		os.Stderr.WriteString("Empty fields inputted in add-book.html.")
		http.Redirect(w, r, "/add-book.html", http.StatusFound)
	} else {
		book.AddBook(bookTitle, bookAuthF, bookAuthL)

		p := LoadPage(members, books)
		RenderTemplate(w, "bookCreated", p)
	}
}

//Handles the edit book page
func EditBookHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	id := r.URL.Path[11:]

	books = book.GetBookById(id)

	if len(books) < 1 {
		http.Redirect(w, r, "/books.html", http.StatusFound)
	} else {
		p := LoadPage(members, books)
		RenderTemplate(w, "edit-book", p)
	}
}

//Handles the edited book page
func BookEditedHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	bookId := helper.HTMLClean(r.FormValue("bookId"))
	bookTitle := helper.HTMLClean(r.FormValue("title"))
	bookAuthF := helper.HTMLClean(r.FormValue("fName"))
	bookAuthL := helper.HTMLClean(r.FormValue("lName"))

	//
	//Consider creating method to check if length == 0
	//
	//Maybe a method that takes an undetermined # of inputs and
	//	returns a false if any are == 0?
	//

	if len(bookId) == 0 || len(bookTitle) == 0 || len(bookAuthF) == 0 || len(bookAuthL) == 0 {
		os.Stderr.WriteString("Empty fields inputted in edit-book.html.")
		http.Redirect(w, r, "/edit-book.html", http.StatusFound)
	} else {
		//
		//Move Library Env Var to a properties file
		//Move all database functions to backing service
		//
		db, err := sql.Open("mysql", os.Getenv("LIBRARY"))
		helper.CheckErr(err)
		defer db.Close()

		stmt, err := db.Prepare("UPDATE books SET book_title = ?, book_authfname = ?, book_authlname = ? WHERE book_id = ?")
		helper.CheckErr(err)

		stmt.Exec(bookTitle, bookAuthF, bookAuthL, bookId)

		p := LoadPage(members, books)
		RenderTemplate(w, "bookEdited", p)
	}
}

//Handles the delete book page
func DeleteBookHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	id := r.URL.Path[13:]

	books = book.GetBookById(id)

	if len(books) < 1 {
		http.Redirect(w, r, "/books.html", http.StatusFound)
	} else {
		p := LoadPage(members, books)
		RenderTemplate(w, "delete-book", p)
	}
}

//Handles the deleted book page
func BookDeletedHandler(w http.ResponseWriter, r *http.Request) {
	id := helper.HTMLClean(r.FormValue("bookId"))

	//
	//Move Database functionality to backing service
	//

	db, err := sql.Open("mysql", os.Getenv("LIBRARY"))
	helper.CheckErr(err)
	defer db.Close()

	//Log transaction
	stmt, err := db.Prepare("DELETE FROM books WHERE book_id = ?")
	helper.CheckErr(err)

	stmt.Exec(id)

	http.Redirect(w, r, "/books.html", http.StatusFound)
}
