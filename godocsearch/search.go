package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// containsKeyword checks if file contains the keyword value provided.
func containsKeyword(filename, keyword string) bool {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error accessing file '%s': %v\n", filename, err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, keyword) {
			return true
		}
	}

	err = scanner.Err()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning file '%s': %v\n", filename, err)
	}

	return false
}

// printResults displays the results of the search for the user.
func printResults(results []string) {
	if len(results) == 0 {
		fmt.Println("No matching files found")
		return
	}

	fmt.Println("The following files matched the keyword given:")
	for _, result := range results {
		fmt.Println(result)
	}
}

// searchFiles walks the file directory based on the SearchSettings.
func searchFiles(settings SearchSettings) []string {
	var results []string

	// TODO: handle case when not recursive
	walker := func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing path '%s': %v\n", path, err)
			return nil
		}

		if !fileInfo.IsDir() {
			if containsKeyword(path, settings.Keyword) {
				results = append(results, path)
			}
		} else {
			subDirSettings := SearchSettings{Keyword: settings.Keyword, Recursive: settings.Recursive, RootDir: path}
			subDirResults := searchFiles(subDirSettings)
			results = append(results, subDirResults...)
		}

		return nil
	}

	err := filepath.Walk(settings.RootDir, walker)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking directory '%s': %v\n", settings.RootDir, err)
	}

	return results
}
