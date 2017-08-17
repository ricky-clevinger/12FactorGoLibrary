package main

//Author: C Neuhardt
//Last Updated: 8/3/2017

import (
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"handlers"
)




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

	http.HandleFunc("/index.html", handlers.MakeHandler(handlers.IndexHandler))
	http.HandleFunc("/search", handlers.MakeHandler(handlers.SearchHandler))
	http.HandleFunc("/admin.html", handlers.MakeHandler(handlers.AdminHandler))
	http.HandleFunc("/books.html", handlers.MakeHandler(handlers.BooksHandler))
	http.HandleFunc("/add-book.html", handlers.MakeHandler(handlers.AddBookHandler))
	http.HandleFunc("/bookCreated", handlers.MakeHandler(handlers.BookCreatedHandler))
	http.HandleFunc("/edit-book/", handlers.MakeHandler(handlers.EditBookHandler))
	http.HandleFunc("/bookEdited", handlers.MakeHandler(handlers.BookEditedHandler))
	http.HandleFunc("/delete-book/", handlers.MakeHandler(handlers.DeleteBookHandler))
	http.HandleFunc("/bookDeleted", handlers.MakeHandler(handlers.BookDeletedHandler))
	http.HandleFunc("/members.html", handlers.MakeHandler(handlers.MembersHandler))
	http.HandleFunc("/add-member.html", handlers.MakeHandler(handlers.AddMemberHandler))
	http.HandleFunc("/memberCreated", handlers.MakeHandler(handlers.MemberCreatedHandler))
	http.HandleFunc("/edit-member/", handlers.MakeHandler(handlers.EditMemberHandler))
	http.HandleFunc("/memberEdited", handlers.MakeHandler(handlers.MemberEditedHandler))
	http.HandleFunc("/delete-member/", handlers.MakeHandler(handlers.DeleteMemberHandler))
	http.HandleFunc("/memberDeleted", handlers.MakeHandler(handlers.MemberDeletedHandler))
	http.HandleFunc("/test.html", handlers.MakeHandler(handlers.TestHandler))
	http.HandleFunc("/checkout.html", handlers.MakeHandler(handlers.CheckoutHandler))
	http.HandleFunc("/checkedout", handlers.MakeHandler(handlers.CheckedoutHandler))
	http.HandleFunc("/checkin.html", handlers.MakeHandler(handlers.CheckinHandler))
	http.HandleFunc("/checkedin", handlers.MakeHandler(handlers.CheckedinHandler))
}

func main() {
	http.HandleFunc("/", redirect)
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	handles()
	http.ListenAndServe(":8080", nil)
}
