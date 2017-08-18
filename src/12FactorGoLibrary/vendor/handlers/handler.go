package handlers

//Created By: Ricky Clevinger
//Updated On: 8/17/2017
//Last Updated By: Ricky Clevinger

import (
	"book"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"helper"
	"html"
	"html/template"
	"member"
	"net/http"
	"os"
	"regexp"
	"time"
)

var validPath = regexp.MustCompile("^/(index.html|search|results.html|admin.html|books.html|add-book.html|bookCreated|edit-book/[0-9]+|bookEdited|delete-book/[0-9]+|bookDeleted|members.html|add-member.html|memberCreated|edit-member/[0-9]+|memberEdited|delete-member/[0-9]+|memberDeleted|test.html|checkout.html|checkedout|checkin.html|checkedin)$")
var templates = template.Must(template.ParseFiles("views/index.html", "views/admin.html", "views/books.html", "views/add-book.html", "views/bookCreated.html", "views/edit-book.html", "views/bookEdited.html", "views/delete-book.html", "views/members.html", "views/add-member.html", "views/memberCreated.html", "views/edit-member.html", "views/memberEdited.html", "views/delete-member.html", "views/test.html", "views/checkout.html", "views/checkedout.html", "views/checkin.html", "views/checkedin.html", "views/results.html"))

type Page struct {
	Members []member.Member
	Books   []book.Book
}

func LoadPage(members []member.Member, books []book.Book) *Page {
	if len(members) > 0 {
		fmt.Println("Loaded member #: ", len(members))
	}
	if len(books) > 0 {
		fmt.Println("Loaded book #: ", len(books))
	}
	return &Page{Members: members, Books: books}
}

//Renders HTML page
func RenderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//Validates path and calls handler
func MakeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
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
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	p := LoadPage(members, books)
	RenderTemplate(w, "index", p)
}

//Handles the admin page
func AdminHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	p := LoadPage(members, books)
	RenderTemplate(w, "admin", p)
}

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
	//bookTitle := html.EscapeString(r.FormValue("title"))
	bookAuthF := helper.HTMLClean(r.FormValue("fName"))
	//bookAuthF := html.EscapeString(r.FormValue("fName"))
	bookAuthL := helper.HTMLClean(r.FormValue("lName"))
	//bookAuthL := html.EscapeString(r.FormValue("lName"))

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

//Handles the members page
func MembersHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	members = member.GetMembers()

	p := LoadPage(members, books)
	RenderTemplate(w, "members", p)
}

//Handles the create member page
func AddMemberHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	p := LoadPage(members, books)
	RenderTemplate(w, "add-member", p)
}

//Handles the create member page
func MemberCreatedHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	memberFName := helper.HTMLClean(r.FormValue("fName"))
	memberLName := helper.HTMLClean(r.FormValue("lName"))

	//
	//Move Database functionality to Backing Service
	//Maybe a method that takes an undetermined # of inputs and
	//      returns a false if any are == 0?
	//

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

		p := LoadPage(members, books)
		RenderTemplate(w, "memberCreated", p)
	}
}

//Handles the edit member page
func EditMemberHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	id := r.URL.Path[13:]

	members = member.GetMemberById(id)

	if len(members) < 1 {
		http.Redirect(w, r, "/members.html", http.StatusFound)
	} else {
		p := LoadPage(members, books)
		RenderTemplate(w, "edit-member", p)
	}
}

//Handles the edited member page
func MemberEditedHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	memberId := helper.HTMLClean(r.FormValue("memId"))
	memberFName := helper.HTMLClean(r.FormValue("fName"))
	memberLName := helper.HTMLClean(r.FormValue("lName"))

	if len(memberId) == 0 || len(memberFName) == 0 || len(memberLName) == 0 {
		os.Stderr.WriteString("Empty fields inputted in edit-member.html.")
		http.Redirect(w, r, "/edit-member.html", http.StatusFound)
	} else {
		db, err := sql.Open("mysql", os.Getenv("LIBRARY"))
		helper.CheckErr(err)
		defer db.Close()

		//Move Database functionality to backing service

		//Log transaction
		stmt, err := db.Prepare("UPDATE member SET member_fname = ?, member_lname = ? WHERE member_id = ?")
		helper.CheckErr(err)

		stmt.Exec(memberFName, memberLName, memberId)

		p := LoadPage(members, books)
		RenderTemplate(w, "memberEdited", p)
	}
}

//Handles the delete member page
func DeleteMemberHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	id := r.URL.Path[15:]

	members = member.GetMemberById(id)

	if len(members) < 1 {
		http.Redirect(w, r, "/members.html", http.StatusFound)
	} else {
		p := LoadPage(members, books)
		RenderTemplate(w, "delete-member", p)
	}
}

//Hanldes the deleted member page
func MemberDeletedHandler(w http.ResponseWriter, r *http.Request) {
	id := html.EscapeString(r.FormValue("memId"))

	//move database functions, change env var location

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
func TestHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	members = member.GetMembers()
	books = book.GetBook()

	p := LoadPage(members, books)
	RenderTemplate(w, "test", p)
}

//Handles the checkout page
func CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	members = member.GetMembers()

	var books []book.Book
	books = book.GetCheckedInBook()

	p := LoadPage(members, books)
	RenderTemplate(w, "checkout", p)
}

//Handles the checkout page
func CheckedoutHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	//
	//Move database functions to the backing service
	//

	current_time := time.Now().Local()

	memberId := helper.HTMLClean(r.FormValue("selPerson"))
	bookId := helper.HTMLClean(r.FormValue("selBook"))
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

	p := LoadPage(members, books)
	RenderTemplate(w, "checkedout", p)
}

//Handles the checkin page
func CheckinHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	members = member.GetMembers()

	var books []book.Book
	books = book.GetCheckedOutBook()

	p := LoadPage(members, books)
	RenderTemplate(w, "checkin", p)
}

//Handles the checkin page
func CheckedinHandler(w http.ResponseWriter, r *http.Request) {
	var members []member.Member
	var books []book.Book

	current_time := time.Now().Local()

	//
	//Move database functions to the backing service
	//

	bookId := helper.HTMLClean(r.FormValue("selBook"))
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

	p := LoadPage(members, books)
	RenderTemplate(w, "checkedin", p)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {

	var books []book.Book
	var members []member.Member
	search := helper.HTMLClean(r.FormValue("s-bar"))

	//
	//Move database functions to the backing service
	//

	if len(search) < 1 {
		os.Stderr.WriteString("Empty fields inputted in home page search.")
		http.Redirect(w, r, "/index.html", http.StatusFound)
	} else {
		books = book.GetSearchedBook(search)
		members = member.GetSearchedMember(search)

		p := LoadPage(members, books)
		RenderTemplate(w, "results", p)
	}
}

//Redirect to index.html
func Redirect(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "/index.html", 301)
}

//Handles
func Handles() {

	http.HandleFunc("/index.html", MakeHandler(IndexHandler))
	http.HandleFunc("/search", MakeHandler(SearchHandler))
	http.HandleFunc("/admin.html", MakeHandler(AdminHandler))
	http.HandleFunc("/books.html", MakeHandler(BooksHandler))
	http.HandleFunc("/add-book.html", MakeHandler(AddBookHandler))
	http.HandleFunc("/bookCreated", MakeHandler(BookCreatedHandler))
	http.HandleFunc("/edit-book/", MakeHandler(EditBookHandler))
	http.HandleFunc("/bookEdited", MakeHandler(BookEditedHandler))
	http.HandleFunc("/delete-book/", MakeHandler(DeleteBookHandler))
	http.HandleFunc("/bookDeleted", MakeHandler(BookDeletedHandler))
	http.HandleFunc("/members.html", MakeHandler(MembersHandler))
	http.HandleFunc("/add-member.html", MakeHandler(AddMemberHandler))
	http.HandleFunc("/memberCreated", MakeHandler(MemberCreatedHandler))
	http.HandleFunc("/edit-member/", MakeHandler(EditMemberHandler))
	http.HandleFunc("/memberEdited", MakeHandler(MemberEditedHandler))
	http.HandleFunc("/delete-member/", MakeHandler(DeleteMemberHandler))
	http.HandleFunc("/memberDeleted", MakeHandler(MemberDeletedHandler))
	http.HandleFunc("/test.html", MakeHandler(TestHandler))
	http.HandleFunc("/checkout.html", MakeHandler(CheckoutHandler))
	http.HandleFunc("/checkedout", MakeHandler(CheckedoutHandler))
	http.HandleFunc("/checkin.html", MakeHandler(CheckinHandler))
	http.HandleFunc("/checkedin", MakeHandler(CheckedinHandler))
}
