package helper

import "testing"

func test(t *testing.T, expected bool, messages ...interface{}) {
	if !expected {
		t.Error(messages)
	}
}

var testData = []struct {
	from string
	to   string
}{
	{"a", "a"},
	{"snake", "snake"},
	{"A", "a"},
	{"ID", "id"},
	{"Snake", "snake"},
	{"SnakeTest", "snake_test"},
	{"SnakeID", "snake_id"},
	{"SnakeIDGoogle", "snake_id_google"},
	{"LinuxFile", "linux_file"},
}

func TestToSnake(t *testing.T) {
	for _, item := range testData {
		result := ToSnake(item.from)
		test(t, result == item.to,
			"Expected result for", item.from, "->", item.to, "got:", result)
	}
}
