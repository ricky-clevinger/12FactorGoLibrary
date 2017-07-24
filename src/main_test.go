package main

import "testing"

func TestLoadPage(t *testing.T) {
	memberIds := []int{1}
	memberNames := []string{"spider man"}

	page := loadPage(memberIds, memberNames)

	if page.MemberIds[0] != 1 && page.MemberNames[0] != "spider man" {
		t.Errorf("Page did not load based on data.")
	}
}
