package model

import "testing"

func TestCase1(t *testing.T) {
	model := Load("file://testdata/case-1.json")
	Equals(t, "Thing", model.AT_Type)
}

func Equals(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Log(expected, " != ", actual)
		t.Fail()
	}
}
