package main

import (
	"flag"
	"fmt"
	"os"
)

// checkKeyword verifies that the keyword has been provided and exits if it is not defined.
func checkKeyword(keyword string) {
	if keyword == "" {
		printUsage()
	}
}

// parseFlags parses the input from the CLI and returns a SearchSettings object.
func parseFlags() SearchSettings {
	var settings SearchSettings

	flag.StringVar(&settings.Keyword, "keyword", "", "Keyword to search for")
	flag.BoolVar(&settings.Recursive, "recursive", false, "Enable recursive search")
	flag.StringVar(&settings.Directory, "dir", ".", "Directory to search (optional)")

	flag.Parse()

	return settings
}

// printUsage displays instructions for the user if the CLI input is incorrect.
func printUsage() {
	fmt.Println("Usage: search -keyword <keyword> [-dir <dir>] [-recursive]")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	settings := parseFlags()
	checkKeyword(settings.Keyword)

	results := searchFiles(settings)
	printResults(results)
}
