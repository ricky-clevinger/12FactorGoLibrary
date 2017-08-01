package book

// Author: Jonathan Quilliams
//Created: 7/25/17
// Edited: 8/1/17
//Purpose: - Query book information/SQL from databse
//		   - Create a Book type
//		   - func NewSlice() -- Creates New slice of Book [NOT USED]
//		   - Getters of Book_id, Book_title, Book_authF, Book_authL, Library_id, Book_check, Mid, and Book_out_date [NOT USED]

import (
	"os"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

//Gets connection string as specified in env vars
var connectionString = os.Getenv("LIBRARY")

//Book Type
type Book struct{
	Book_id int
	Book_title string
	Book_authF string
	Book_authL string
	Library_id int
	Book_check int
	Mid int
	Book_out_date sql.NullString //Data pulled from db is either a date-string or null.
}

//Create new Book Slice Type
//Typically used outside of book.go
//e.g. BookVar := book.NewSlice()
func NewSlice() *[]Book {return new([]Book)}

//Get Book
//Returns Book Slice of multiple books
func GetBook() []Book {
	
	var books []Book //Holds Slice of Book Type

	//DB Connection
	db, err := sql.Open("mysql", connectionString)
	checkErr(err)
	defer db.Close() //Close after func GetBook ends

	//Grab entire rows of data within a query
	bookRows, err := db.Query("SELECT Book_Id, Book_Title, Book_AuthFName, Book_AuthLName, Library_Id, Book_Check, Mid, date_format(Book_Out_Date, '%Y-%m-%d') FROM books")
	//Check for Errors in DB Query
	checkErr(err)
	//For Every Book Row that's not null/nil place
	for bookRows.Next() {
		b := Book{} //book type
		err = bookRows.Scan(&b.Book_id, &b.Book_title, &b.Book_authF, &b.Book_authL, &b.Library_id, &b.Book_check, &b.Mid, &b.Book_out_date)
		checkErr(err)
		if b.Book_out_date.Valid{
			books = append(books, b)
		} else {
			b.Book_out_date.String = ""
			books = append(books, b)
		}

	}

	return books
}// end GetBook()

/* //Get Book Id
func GetId() []int {
	
	var bookIds []int

	db, err := sql.Open("mysql", connectionString)
	checkErr(err)
	defer db.Close()
	
	bookIdRows, err := db.Query("SELECT Book_Id FROM books")
	checkErr(err)
	for bookIdRows.Next() {
		var hBookId int
		err = bookIdRows.Scan(&hBookId)
		checkErr(err)
		bookIds = append(bookIds, hBookId)
	}
	
	return bookIds
}

//Get Book Title
func GetTitle() []string {
	
	var bookTitles []string

	db, err := sql.Open("mysql", connectionString)
	checkErr(err)
	defer db.Close()
	
	bookTitleRows, err := db.Query("SELECT Book_Title FROM books")
	checkErr(err)
	for bookTitleRows.Next() {
		var hBookTitle string
		err = bookTitleRows.Scan(&hBookTitle)
		checkErr(err)
		bookTitles = append(bookTitles, hBookTitle)
	}
	
	return bookTitles
}

//Get Book Author's First Name
func GetFName() []string {
	
	var bookFNames []string

	db, err := sql.Open("mysql", connectionString)
	checkErr(err)
	defer db.Close()
	
	bookFNameRows, err := db.Query("SELECT Book_AuthFName FROM books")
	checkErr(err)
	for bookFNameRows.Next() {
		var hBookFName string
		err = bookFNameRows.Scan(&hBookFName)
		checkErr(err)
		bookFNames = append(bookFNames, hBookFName)
	}
	
	return bookFNames
}

//Get Book Author's Last Name
func GetLName() []string {
	
	var bookLNames []string

	db, err := sql.Open("mysql", connectionString)
	checkErr(err)
	defer db.Close()
	
	bookLNameRows, err := db.Query("SELECT Book_AuthLName FROM books")
	checkErr(err)
	for bookLNameRows.Next() {
		var hBookLName string
		err = bookLNameRows.Scan(&hBookLName)
		checkErr(err)
		bookLNames = append(bookLNames, hBookLName)
	}
	
	return bookLNames
}

//Get Book Library Id
func GetLibId() []int {
	
	var bookLibIds []int

	db, err := sql.Open("mysql", connectionString)
	checkErr(err)
	defer db.Close()
	
	bookLibIdRows, err := db.Query("SELECT Library_Id FROM books")
	checkErr(err)
	for bookLibIdRows.Next() {
		var hBookLibId int
		err = bookLibIdRows.Scan(&hBookLibId)
		checkErr(err)
		bookLibIds = append(bookLibIds, hBookLibId)
	}
	
	return bookLibIds
}

//Get Book's Checked-Out Status
func GetBookCheck() []int {
	
	var bookChecks []int

	db, err := sql.Open("mysql", connectionString)
	checkErr(err)
	defer db.Close()
	
	bookCheckRows, err := db.Query("SELECT Book_Check FROM books")
	checkErr(err)
	for bookCheckRows.Next() {
		var hBookCheck int
		err = bookCheckRows.Scan(&hBookCheck)
		checkErr(err)
		bookChecks = append(bookChecks, hBookCheck)
	}
	
	return bookChecks
}

//Get Book's Member Id
func GetMId() []int {
	
	var bookMIds []int

	db, err := sql.Open("mysql", connectionString)
	checkErr(err)
	defer db.Close()
	
	bookMIdRows, err := db.Query("SELECT Mid FROM books")
	checkErr(err)
	for bookMIdRows.Next() {
		var hBookMId int
		err = bookMIdRows.Scan(&hBookMId)
		checkErr(err)
		bookMIds = append(bookMIds, hBookMId)
	}
	
	return bookMIds
}

//Get Book's Checked-Out Date
func GetCheckOutDate() []string {
	
	var bookDates []string

	db, err := sql.Open("mysql", connectionString)
	checkErr(err)
	defer db.Close()
	
	bookDateRows, err := db.Query("SELECT CAST(Book_Out_Date AS CHAR) FROM books")
	checkErr(err)
	for bookDateRows.Next() {
		var hBookDate string
		err = bookDateRows.Scan(&hBookDate)
		checkErr(err)
		bookDates = append(bookDates, hBookDate)
	}
	
	return bookDates
}

//Get Book Titles that aren't Checked-Out
func GetBooksNotOut() []string {
	
	var booksNotOut []string

	db, err := sql.Open("mysql", connectionString)
	checkErr(err)
	defer db.Close()
	
	bookNotOutRows, err := db.Query("SELECT Book_Title FROM books WHERE Book_Check = 1")
	checkErr(err)
	for bookCheckRows.Next() {
		var hBookCheck string
		err = bookCheckRows.Scan(&hBookCheck)
		checkErr(err)
		booksNotOut = append(booksNotOut, hBookCheck)
	}
	
	return booksNotOut
} */


//Checks for errors
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}