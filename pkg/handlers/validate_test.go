package handlers

import "testing"

func test(t *testing.T, expected bool, messages ...interface{}) {
	if !expected {
		t.Error(messages)
	}
}

var testData = []struct {
	str    string
	result bool
}{
	{"a", true},
	{"test123", true},
	{"test-123", true},
	{"test_456", true},
	{"old.value-23", true},
	{"12345", true},
	{"#$%&", false},
}

func TestIsValidString(t *testing.T) {
	for _, item := range testData {
		result := isValidString(item.str)
		test(t, result == item.result,
			"Expected result for", item.str, "->", item.result, "got:", result)
	}
}
