package utils

import (
	"slices"
	"testing"
)

func TestFullmatch(t *testing.T) {
	if !FullMatch(`\d+`, "124345") {
		t.Errorf("Test Fullmatch failed!")
	}
	if FullMatch(`\d+`, " 124345") {
		t.Errorf("Test Fullmatch failed!")
	}
}

func TestMatch(t *testing.T) {
	if !Match(`\d+`, "124345") {
		t.Errorf("Test Match failed!")
	}
	if !Match(`\d+`, "Hello124345 ") {
		t.Errorf("Test Match failed!")
	}
}

func TestExtract(t *testing.T) {
	if Extract(`\d+`, "where124345hello") != "124345" {
		t.Errorf("Test Extract failed!")
	}
	s := FindAll(`\d+`, "hello124345world123ok")
	if !slices.Equal(s, []string{"124345", "123"}) {
		t.Errorf("Test FindAll failed!,%s", s)
	}
}
