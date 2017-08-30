package handlers

//Created By: Ricky Clevinger
//Updated On: 8/17/2017
//Last Updated By: Ricky Clevinger

import (
	"book"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"helper"
	"member"
	"net/http"
	"os"
	"html"
	"crypto/sha1"
        "encoding/base64"
)


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
	memberEmail := helper.HTMLClean(r.FormValue("email"))
	memberPass := helper.HTMLClean(r.FormValue("password"))

	//
	//Move Database functionality to Backing Service
	//Maybe a method that takes an undetermined # of inputs and
	//      returns a false if any are == 0?
	//

	if len(memberFName) == 0 || len(memberLName) == 0 {
		os.Stderr.WriteString("Empty fields inputted in add-member.html.")
		http.Redirect(w, r, "/register.html", http.StatusFound)
	} else {
		db, err := sql.Open("mysql", os.Getenv("LIBRARY"))
		helper.CheckErr(err)
		defer db.Close()

		stmt, err := db.Prepare("INSERT INTO member (member_fname, member_lname, Email, Password) VALUES (?, ?, ?, ?)")
		helper.CheckErr(err)

		bv := []byte(memberPass) 
		hasher := sha1.New()
    		hasher.Write(bv)

		stmt.Exec(memberFName, memberLName, memberEmail, base64.URLEncoding.EncodeToString(hasher.Sum(nil)))

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
