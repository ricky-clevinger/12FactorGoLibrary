package main

import (
	"testing"
	"member"
	"book"
	"fmt"
)

//Begin main package testing
func TestLoadPage(t *testing.T) {
	members := member.GetMembers()
	books := book.GetBook()

	page := loadPage(members, books)

	if len(page.Members) < 1 || len(page.Books) < 1 {
		t.Errorf("Page did not load based on data.")
	} else {
		fmt.Println("index.go Func loadPage PASS")
	}
}
//End main package testing

//Begin member package testing
func TestGetMembers(t *testing.T) {
	members := member.GetMembers()

	if len(members) < 1 {
		t.Errorf("No members loaded from database.")
	} else {
		fmt.Println("member.go Func GetMembers PASS")
	}
}

func TestGetSearchedMember(t *testing.T) {
	members := member.GetSearchedMember("Goku")

	if len(members) < 1 {
		t.Errorf("No members loaded from database.")
	} else {
		fmt.Println("member.go Func GetSearchMember PASS")
	}
}
//End member package testing.

//Begin book package testing
func TestGetBook(t *testing.T) {
	books := book.GetBook()

	if len(books) < 1 {
		t.Errorf("No books loaded from database.")
	} else {
		fmt.Println("book.go Func GetBook PASS")
	}
}

func TestGetSearchBook(t *testing.T) {
	books := book.GetSearchedBook("Harry")

	if len(books) < 1 {
		t.Errorf("No books loaded from database.")
	} else {
		fmt.Println("book.go Func GetSearchedBook PASS")
	}
}
//End book package testing