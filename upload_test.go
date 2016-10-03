package main

import (
	"strings"
	"testing"
)

func TestJsonAnswer(t *testing.T) {
	if string(jsonAnswer(map[string]string{"Test": "Das ist ein Test"})) !=
		"{\"Test\":\"Das ist ein Test\"}" {
		t.Error("Uncorrect json")
	}
}

func TestIsHash(t *testing.T) {
	hash := "e9b4a070e1be7dbec8b340ef80744d32f8d3cb9a9d1f89fed225037b9eaf0a271876adadc7485a6090aa2e8c0e30984c26710a501ce889cccb363d1cb28f087b"
	if !isHash(hash, strings.NewReader("Das ist ein Test")) {
		t.Error("Correct hash, but isHash returns false")
	}
	if isHash(hash, strings.NewReader("Das ist noch ein Test")) {
		t.Error("Incorrect hash, but isHash returns true")
	}
}
