package book

// Author: Jonathan Quilliams
//Created: 7/25/17
// Edited: 8/3/17
//Purpose: - Query book information/SQL from database
//		   - Create a Book type
//		   - func NewSlice() -- Creates New slice of Book [NOT USED]

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"os"
	"strconv"
	"helper"
)

//Gets connection string as specified in env vars

//Move the assignment of the env var

var connectionString = os.Getenv("LIBRARY")

//Connects to RESTful API

//move this into a properties file and set as a variable here
var apiUrl = "http://localhost:8081/books"

//Book Type
type Book struct {
	Book_id       int            `json:"Book_id"`
	Book_title    string         `json:"Book_title"`
	Book_authF    string         `json:"Book_authF"`
	Book_authL    string         `json:"Book_authL"`
	Library_id    int            `json:"Library_id"`
	Book_check    string         `json:"Book_check"`
	Mid           int            `json:"Mid"`
	Book_out_date sql.NullString `json:"Book_out_date"`
	Member_fname  sql.NullString `json:"Member_fname"`
	Member_lname  sql.NullString `json:"Member_lname"`
}

//Create new Book Slice Type
//Typically used outside of book.go
//e.g. BookVar := book.NewSlice()
func NewSlice() *[]Book { return new([]Book) }

//Get Books
func GetBook() []Book {
	var books []Book

	request, err := http.NewRequest("GET", apiUrl, nil)
	helper.CheckErr(err)

	client := &http.Client{}

	response, err := client.Do(request)
	helper.CheckErr(err)
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&books)
	helper.CheckErr(err)

	return books
}

//Get Book
//Returns Book Slice of multiple books

// ---> Chris moved this to the backing service

/*func GetBook() []Book {

	var books []Book //Holds Slice of Book Type

	//DB Connection
	db, err := sql.Open("mysql", connectionString)
	checkErr(err)
	defer db.Close() //Close after func GetBook ends

	//Grab entire rows of data within a query
	bookRows, err := db.Query("SELECT b.Book_Id, b.Book_Title, b.Book_AuthFName, b.Book_AuthLName, b.Library_Id, CASE WHEN b.Book_check = 1 THEN 'In' ELSE 'Out' END AS B_check, b.Mid, date_format(b.Book_Out_Date, '%Y-%m-%d'), m.Member_fname, m.Member_lname FROM books b LEFT JOIN member m ON b.Mid = m.Member_id")
	//Check for Errors in DB Query
	checkErr(err)
	//For Every Book Row that's not null/nil place
	for bookRows.Next() {
		b := Book{} //book type
		err = bookRows.Scan(&b.Book_id, &b.Book_title, &b.Book_authF, &b.Book_authL, &b.Library_id, &b.Book_check, &b.Mid, &b.Book_out_date, &b.Member_fname, &b.Member_lname)
		checkErr(err)
		if b.Book_out_date.Valid{
			books = append(books, b)
		} else {
			b.Book_out_date.String = ""
			books = append(books, b)
		}

	}

	return books
}// end GetBook()*/

//Get Member by ID
func GetBookById(id string) []Book {
	var book []Book
	var books []Book
	intId, err := strconv.Atoi(id)
	helper.CheckErr(err)

	request, err := http.NewRequest("GET", apiUrl, nil)
	helper.CheckErr(err)

	client := &http.Client{}

	response, err := client.Do(request)
	helper.CheckErr(err)
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&books)
	helper.CheckErr(err)

	for i := 0; i < len(books); i++ {
		booksId := books[i].Book_id

		if booksId == intId {
			book = append(book, books[i])
		}
	}

	return book
}

// -----> chris moved this to the backing service

/*func GetBookById(id string) []Book {

	var book []Book //Holds Slice of Book Type

	//DB Connection
	db, err := sql.Open("mysql", connectionString)
	checkErr(err)
	defer db.Close() //Close after func GetBook ends

	//Grab entire rows of data within a query
	bookRows, err := db.Query("SELECT b.Book_Id, b.Book_Title, b.Book_AuthFName, b.Book_AuthLName, b.Library_Id, CASE WHEN b.Book_check = 1 THEN 'In' ELSE 'Out' END AS B_check, b.Mid, date_format(b.Book_Out_Date, '%Y-%m-%d'), m.Member_fname, m.Member_lname FROM books b LEFT JOIN member m ON b.Mid = m.Member_id WHERE Book_Id = ?", id)
	//Check for Errors in DB Query
	checkErr(err)
	//For Every Book Row that's not null/nil place
	for bookRows.Next() {
		b := Book{} //book type
		err = bookRows.Scan(&b.Book_id, &b.Book_title, &b.Book_authF, &b.Book_authL, &b.Library_id, &b.Book_check, &b.Mid, &b.Book_out_date, &b.Member_fname, &b.Member_lname)
		checkErr(err)
		if b.Book_out_date.Valid{
			book = append(book, b)
		} else {
			b.Book_out_date.String = ""
			book = append(book, b)
		}

	}

	return book
}// end GetBook()*/

//Returns Book Slice of multiple books that are not checked out
func GetCheckedInBook() []Book {
	var checkedBooks []Book
	var books []Book

	request, err := http.NewRequest("GET", apiUrl, nil)
	helper.CheckErr(err)

	client := &http.Client{}

	response, err := client.Do(request)
	helper.CheckErr(err)
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&books)
	helper.CheckErr(err)

	for i := 0; i < len(books); i++ {
		if books[i].Book_check == "In" {
			checkedBooks = append(checkedBooks, books[i])
		}
	}

	return checkedBooks
}

//Returns Book Slice of multiple books that are checked out
func GetCheckedOutBook() []Book {
	var checkedBooks []Book
	var books []Book

	request, err := http.NewRequest("GET", apiUrl, nil)
	helper.CheckErr(err)

	client := &http.Client{}

	response, err := client.Do(request)
	helper.CheckErr(err)
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&books)
	helper.CheckErr(err)

	for i := 0; i < len(books); i++ {
		if books[i].Book_check == "Out" {
			checkedBooks = append(checkedBooks, books[i])
		}
	}

	return checkedBooks
}

// Return Books under search params

//move this to the backing service

func GetSearchedBook(s string) []Book {

	var books []Book //Hold Slice of Book Type
	search := fmt.Sprintf("%s%s%s", "%", s, "%")

	//DB Connection
	db, err := sql.Open("mysql", connectionString)
	helper.CheckErr(err)
	defer db.Close() //Close after func GetBook ends

	//Prepare entire rows of data within a query
	bookRows, err := db.Query("SELECT b.Book_Id, b.Book_Title, b.Book_AuthFName, b.Book_AuthLName, b.Library_Id, CASE WHEN b.Book_check = 1 THEN 'In' ELSE 'Out' END AS B_check, b.Mid, date_format(b.Book_Out_Date, '%Y-%m-%d'), m.Member_fname, m.Member_lname FROM books b LEFT JOIN member m ON b.Mid = m.Member_id WHERE book_title like ?", search)

	//Check for Errors in DB the Query
	helper.CheckErr(err)

	//For Every Book Row that's not null/nil place
	for bookRows.Next() {
		b := Book{} //book type
		err = bookRows.Scan(&b.Book_id, &b.Book_title, &b.Book_authF, &b.Book_authL, &b.Library_id, &b.Book_check, &b.Mid, &b.Book_out_date, &b.Member_fname, &b.Member_lname)
		helper.CheckErr(err)
		if b.Book_out_date.Valid {
			books = append(books, b)
		} else {
			b.Book_out_date.String = ""
			books = append(books, b)
		}

	}

	return books
}

//INSERT New Row into books TABLE

//move this to the backing service

func AddBook(title string, authF string, authL string) {
	Book_title := title
	Book_authF := authF
	Book_authL := authL

	db, err := sql.Open("mysql", os.Getenv("LIBRARY"))
	helper.CheckErr(err)
	defer db.Close()

	//Insert new book instance into table
	stmt, err := db.Prepare("INSERT INTO books (Book_Title, Book_AuthFName, Book_AuthLName, Library_Id, Book_Check, Mid) VALUES (?, ?, ?, 1, 1, 0)")
	helper.CheckErr(err)

	stmt.Exec(Book_title, Book_authF, Book_authL)
}

//UPDATE Row in books TABLE
/*func EditBook(bookId string, title string, authF string, authL string) {
	Book_id := bookId
	Book_title := title
	Book_authF := authF
	Book_authL := authL

	db, err := sql.Open("mysql", os.Getenv("LIBRARY"))
	checkErr(err)
	defer db.Close()

	//Insert new book instance into table
	stmt, err := db.Prepare("UPDATE books SET Book_Title = ? Book_AuthFName = ? Book_AuthLName = ? WHERE Book_Id = ?")
	checkErr(err)

	stmt.Exec(Book_title, Book_authF, Book_authL, Book_id)
}*/

/*//Checks for errors
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}*/
