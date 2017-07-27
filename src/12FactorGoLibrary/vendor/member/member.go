package member

//Author: C Neuhardt
//Last Updated: 7/26/2017

import (
	"os"
	"fmt"	
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

//Gets the connection string "cgidevlib:Password1@tcp(cgiprojdevlibrary.cxyeb3wmov3g.us-east-1.rds.amazonaws.com:9871)/cgiprojdevlibrary"
var connectionString = os.Getenv("LIBRARY")

//Return an array of IDs for each member in the db
func GetIds() []int {
	var memberIds []int
	
	db, err := sql.Open("mysql", connectionString)
	checkErr(err)
	defer db.Close()
	
	memberIdRows, err := db.Query("SELECT member_id FROM member")
	checkErr(err)
	for memberIdRows.Next() {
		var MemberId int
		err = memberIdRows.Scan(&MemberId)
		checkErr(err)
		memberIds = append(memberIds, MemberId)
	}

	return memberIds
}

//Return an array of First Names for each member in the db
func GetFNames() []string {
	var memberFNames []string
	
	db, err := sql.Open("mysql", connectionString)
	checkErr(err)
	defer db.Close()
	
	memberIdRows, err := db.Query("SELECT member_fname FROM member")
	checkErr(err)
	for memberIdRows.Next() {
		var Member_FName string
		err = memberIdRows.Scan(&Member_FName)
		checkErr(err)
		memberFNames = append(memberFNames, Member_FName)
	}

	return memberFNames
}

//Return an array of Last Names for each member in the db
func GetLNames() []string {
	var memberLNames []string
	
	db, err := sql.Open("mysql", connectionString)
	checkErr(err)
	defer db.Close()
	
	memberIdRows, err := db.Query("SELECT member_lname FROM member")
	checkErr(err)
	for memberIdRows.Next() {
		var Member_LName string
		err = memberIdRows.Scan(&Member_LName)
		checkErr(err)
		memberLNames = append(memberLNames, Member_LName)
	}

	return memberLNames
}

//Return First Name for a member using member_id
func GetFNameById(memberId int) string {
	db, err := sql.Open("mysql", connectionString)
	checkErr(err)
	defer db.Close()

	queryString := fmt.Sprintf("SELECT member_fname FROM member WHERE member_id = %d", memberId)
	memberRows, err := db.Query(queryString)
	checkErr(err)
	for memberRows.Next() {
		var Member_FName string
		err = memberRows.Scan(&Member_FName)
		checkErr(err)
		return Member_FName
	}

	return ""
}

//Return Last Name for a member using member_id
func GetLNameById(memberId int) string {
	db, err := sql.Open("mysql", connectionString)
	checkErr(err)
	defer db.Close()

	queryString := fmt.Sprintf("SELECT member_lname FROM member WHERE member_id = %d", memberId)
	memberRows, err := db.Query(queryString)
	checkErr(err)
	for memberRows.Next() {
		var Member_LName string
		err = memberRows.Scan(&Member_LName)
		checkErr(err)
		return Member_LName
	}

	return ""
}

//Checks for errors
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
