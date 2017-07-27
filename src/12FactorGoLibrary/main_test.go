package main

import (
	"testing"
	"member"
)

//Begin main package testing
func TestLoadPage(t *testing.T) {
	memberIds := []int{1}
	memberFNames := []string{"spider man"}

	page := loadPage(memberIds, memberFNames)

	if page.MemberIds[0] != 1 && page.MemberFNames[0] != "spider man" {
		t.Errorf("Page did not load based on data.")
	}
}
//End main package testing

//Begin member package testing
func TestGetIds(t *testing.T) {
	memberIds := member.GetIds()

	if len(memberIds) < 1 {
		t.Errorf("No member IDs loaded from database.")
	}
}

func TestGetFNames(t *testing.T) {
	memberFNames := member.GetFNames()

	if len(memberFNames) < 1 {
		t.Errorf("No first names loaded from database.")
	}
}

func TestGetLNames(t *testing.T) {
	memberLNames := member.GetLNames()

	if len(memberLNames) < 1 {
		t.Errorf("No last names loaded from database.")
	}
}

func TestGetFNameById(t *testing.T) {
	FName := member.GetFNameById(1)

	if FName != "Ricky" {
		t.Errorf("Incorrect First Name retrieved when searching for member with ID = 1")
	}
}

func TestGetLNameById(t *testing.T) {
	LName := member.GetLNameById(1)

	if LName != "Clevinger" {
		t.Errorf("Incorrect Last Name retrieved when searching for member with ID = 1")
	}
}
//End member package testing.