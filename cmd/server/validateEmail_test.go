package main

import "testing"

func TestValidateEmail(t *testing.T) {
	test_case := "normal@gmail.com"
	if !validateEmail(test_case) {
		t.Errorf("Email '%s' tested invalid, should be valid", test_case)
	}

	test_case = "@"
	if validateEmail(test_case) {
		t.Error("Email with just the @ sign should not get validated")
	}

	test_case = "onlyleft@"
	if validateEmail(test_case) {
		t.Error("Email with only the left side of the @ should not get validated")
	}

	test_case = "@onlyright"
	if validateEmail(test_case) {
		t.Error("Email with only the right side of the @ should not get validated")
	}

	test_case = "not an email address"
	if validateEmail(test_case) {
		t.Error("An email address with no @ should not get validated")
	}
}
