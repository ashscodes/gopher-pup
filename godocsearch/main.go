package main

import (
	"flag"
	"fmt"
	"os"
)

// SearchSettings contains the parsed CLI flag values.
type SearchSettings struct {
	Keyword   string
	Recursive bool
	RootDir   string
}

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
	flag.StringVar(&settings.RootDir, "root", ".", "Root directory to search (optional)")

	flag.Parse()

	return settings
}

// printUsage displays instructions for the user if the CLI input is incorrect.
func printUsage() {
	fmt.Println("Usage: search -keyword <keyword> [-root <rootDir>] [-recursive]")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	settings := parseFlags()
	checkKeyword(settings.Keyword)

	results := searchFiles(settings)
	printResults(results)
}
