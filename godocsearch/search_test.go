package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Uncomment and import the main package if tests are in a different directory
// import "./"

func BenchmarkContainsKeyword(b *testing.B) {
	filename := "testdata/test1.txt"
	keyword := "hello"

	for i := 0; i < b.N; i++ {
		containsKeyword(filename, keyword)
	}
}
func BenchmarkSearchFiles(b *testing.B) {
	settings := SearchSettings{
		Keyword:   "hello",
		Recursive: false,
		Directory: "testdata",
	}

	for i := 0; i < b.N; i++ {
		searchFiles(settings)
	}
}

// TestContainsKeyword tests the function containsKeyword
func TestContainsKeyword(t *testing.T) {
	tests := []struct {
		Filename string
		Keyword  string
		Result   bool
	}{
		{"testdata/test1.txt", "hello", true},
		{"testdata/subdir/test2.txt", "world", true},
		{"testdata/test3.txt", "foo", false},
	}

	for _, tc := range tests {
		t.Run(tc.Filename, func(t *testing.T) {
			actual := containsKeyword(tc.Filename, tc.Keyword)
			assert.Equal(t, tc.Result, actual, "Mismatch in keyword presence for file: "+tc.Filename)
		})
	}
}

// TestSearchFiles tests the function searchFiles
func TestSearchFiles(t *testing.T) {
	tests := []struct {
		Name     string
		Result   []string
		Settings SearchSettings
	}{
		{
			"files_containing_hello_when_non-recursive",
			[]string{"testdata/test1.txt"},
			SearchSettings{Keyword: "hello", Recursive: false, Directory: "testdata"},
		},
		{
			"files_containing_hello_when_recursive",
			[]string{"testdata/test1.txt"},
			SearchSettings{Keyword: "hello", Recursive: true, Directory: "testdata"},
		},
		{
			"files_containing_file_when_non-recursive",
			[]string{"testdata/test1.txt", "testdata/test3.txt"},
			SearchSettings{Keyword: "file", Recursive: false, Directory: "testdata"},
		},
		{
			"files_containing_file_when_recursive",
			[]string{"testdata/test1.txt", "testdata/subdir/test2.txt", "testdata/test3.txt", "testdata/subdir/furthersubdir/test4.txt"},
			SearchSettings{Keyword: "file", Recursive: true, Directory: "testdata"},
		},
		{
			"files_containing_This_when_non-recursive",
			[]string{"testdata/test1.txt", "testdata/test3.txt"},
			SearchSettings{Keyword: "This", Recursive: false, Directory: "testdata"},
		},
		{
			"files_containing_This_when_recursive",
			[]string{"testdata/test1.txt", "testdata/subdir/test2.txt", "testdata/test3.txt"},
			SearchSettings{Keyword: "This", Recursive: true, Directory: "testdata"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			actual := searchFiles(tc.Settings)
			assert.ElementsMatch(t, tc.Result, actual, "Mismatch in search results")
		})
	}
}
