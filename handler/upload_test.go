package handler

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
