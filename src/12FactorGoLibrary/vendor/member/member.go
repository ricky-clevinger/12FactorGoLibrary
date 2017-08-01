package member

//Author: C Neuhardt
//Last Updated: 8/1/2017

import (
	"os"
	//"fmt"	
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

//Gets the connection string
var connectionString = os.Getenv("LIBRARY")

//Member Type
type Member struct {
	Member_id int
	Member_fname string
	Member_lname string
}

//Get Members
func GetMembers() []Member {
	
	var members []Member

	db, err := sql.Open("mysql", connectionString)
	checkErr(err)
	defer db.Close()

	memberRows, err := db.Query("SELECT member_id, member_fname, member_lname FROM member")
	checkErr(err)
	
	for memberRows.Next() {
		m := Member{}
		err = memberRows.Scan(&m.Member_id, &m.Member_fname, &m.Member_lname)
		checkErr(err)
		members = append(members, m)
	}

	return members
}

//Checks for errors
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
