package member

//Author: C Neuhardt
//Last Updated: 8/3/2017

import (
	"encoding/json"
	"os"
	"fmt"
	"strconv"
	"database/sql"
	"net/http"
	_ "webPackages/github.com/go-sql-driver/mysql"
	"helper"
	"crypto/sha1"
    	"encoding/base64"
)

//Gets the connection string
var connectionString = os.Getenv("LIBRARY")

//Connects to RESTful API
var url = "http://localhost:8081/members"

//Member Type
type Member struct {
	Member_id int `json:"Member_id"`
	Member_fname string `json:"Member_fname"`
	Member_lname string `json:"Member_lname"`
	Role string `json:"Role"`
}

//Get Members
func GetMembers() []Member {
	var members []Member

	request, err := http.NewRequest("GET", url, nil)
	helper.CheckErr(err)

	client := &http.Client{}

	response, err := client.Do(request)
	helper.CheckErr(err)
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&members)
	helper.CheckErr(err)

	return members
}

//Get Members
func MemberExist(mail, pass string) []Member {
	var members []Member

	bv := []byte(pass) 
	hasher := sha1.New()
    	hasher.Write(bv)

	request, err := http.NewRequest("GET", url + "/login/" + mail + "/" + base64.URLEncoding.EncodeToString(hasher.Sum(nil)), nil)
	helper.CheckErr(err)

	client := &http.Client{}

	response, err := client.Do(request)
	helper.CheckErr(err)
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&members)
	helper.CheckErr(err)

	return members
}

//Get Member by ID
func GetMemberById(id string) []Member {
	var member []Member
	var members []Member
	intId, err := strconv.Atoi(id)
	helper.CheckErr(err)

	request, err := http.NewRequest("GET", url, nil)
	helper.CheckErr(err)

	client := &http.Client{}

	response, err := client.Do(request)
	helper.CheckErr(err)
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&members)
	helper.CheckErr(err)

	for i := 0; i < len(members); i++ {
		membersId := members[i].Member_id

		if membersId == intId {
			member = append(member, members[i])
		}
	}

	return member
}

/*func GetMemberById(id string) []Member {
	
	var member []Member

	db, err := sql.Open("mysql", connectionString)
	checkErr(err)
	defer db.Close()

	memberRows, err := db.Query("SELECT member_id, member_fname, member_lname FROM member WHERE member_id = ?", id)
	checkErr(err)
	
	for memberRows.Next() {
		m := Member{}
		err = memberRows.Scan(&m.Member_id, &m.Member_fname, &m.Member_lname)
		checkErr(err)
		member = append(member, m)
	}

	return member
}*/

//Get Members using search
func GetSearchedMember(s string) []Member {
	
	var members []Member //Hold Slice of Member Type
	search :=  fmt.Sprintf("%s%s%s", "%", s, "%")

	//DB Connection
	db, err := sql.Open("mysql", connectionString)
	helper.CheckErr(err)
	defer db.Close() //Close after func GetSearchedMember(s string) ends

	//Prepare entire rows of data within a query
	memberRows, err := db.Query("SELECT member_id, member_fname, member_lname FROM member WHERE member_fname like ? OR member_lname like ?", search, search)
	
	//Check for Errors in DB the Query
	helper.CheckErr(err)

	//For Every Member Row that's not null/nil place
	for memberRows.Next() {
		m := Member{}
		err = memberRows.Scan(&m.Member_id, &m.Member_fname, &m.Member_lname)
		helper.CheckErr(err)
		members = append(members, m)
	}

	return members
}

/*
//Checks for errors
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
*/
