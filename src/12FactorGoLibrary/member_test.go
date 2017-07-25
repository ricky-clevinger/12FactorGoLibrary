package member

import "testing"

func TestGetIds(t *testing.T) {
	memberIds := getIds()

	if len(memberIds) < 1 {
		t.Errorf("No member IDs loaded from database.")
	}
}

func TestGetFNames(t *testing.T) {
	memberFNames := getFNames()

	if len(memberFNames) < 1 {
		t.Errorf("No first names loaded from database.")
	}
}

func TestGetLNames(t *testing.T) {
	memberLNames := getLNames()

	if len(memberLNames) < 1 {
		t.Errorf("No last names loaded from database.")
	}
}

func TestGetFNameById(t *testing.T) {
	FName := getFNameById(1)

	if FName != "Ricky" {
		t.Errorf("Incorrect First Name retrieved when searching for member with ID = 1")
	}
}

func TestGetLNameById(t *testing.T) {
	LName := getLNameById(1)

	if LName != "Clevinger" {
		t.Errorf("Incorrect Last Name retrieved when searching for member with ID = 1")
	}
}