/*
 * Glimpse - Search Engine Module
 * Author: Pratyush Yadav <pratyushyadav0106@gmail.com>
 * Description: Core search functionality with file traversal and pattern matching
 * License: MIT License
 */
package search

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Result struct {
	File    string
	Line    int
	Content string
}

type SearchOptions struct {
	CaseInsensitive bool
	MaxResults      int
}

func SearchFiles(pattern string, dir string, options SearchOptions) []Result {
	var results []Result
	var mutex sync.Mutex

	maxResults := options.MaxResults
	if maxResults == 0 {
		maxResults = 50
	}

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			name := info.Name()
			if name == ".git" || name == "node_modules" || name == ".vscode" ||
				name == "target" || name == "build" || name == "dist" {
				return filepath.SkipDir
			}
			return nil
		}

		mutex.Lock()
		if len(results) >= maxResults {
			mutex.Unlock()
			return filepath.SkipDir
		}
		mutex.Unlock()

		if isBinaryFile(path) {
			return nil
		}

		fileResults := searchFile(pattern, path, options)

		mutex.Lock()
		results = append(results, fileResults...)
		if len(results) > maxResults {
			results = results[:maxResults]
		}
		mutex.Unlock()

		return nil
	})

	return results
}

func searchFile(pattern string, filename string, options SearchOptions) []Result {
	var results []Result

	file, err := os.Open(filename)
	if err != nil {
		return results
	}
	defer file.Close()

	searchPattern := pattern
	if options.CaseInsensitive {
		searchPattern = strings.ToLower(pattern)
	}

	scanner := bufio.NewScanner(file)
	lineNumber := 1

	for scanner.Scan() {
		if len(results) >= 5 {
			break
		}

		line := scanner.Text()

		searchLine := line
		if options.CaseInsensitive {
			searchLine = strings.ToLower(line)
		}

		if strings.Contains(searchLine, searchPattern) {
			results = append(results, Result{
				File:    filename,
				Line:    lineNumber,
				Content: line,
			})
		}
		lineNumber++
	}

	return results
}

func isBinaryFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	binaryExts := []string{
		".exe", ".dll", ".so", ".dylib", ".a", ".o",
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".ico",
		".pdf", ".zip", ".tar", ".gz", ".7z",
		".mp3", ".mp4", ".avi", ".mov",
	}

	for _, binaryExt := range binaryExts {
		if ext == binaryExt {
			return true
		}
	}

	file, err := os.Open(filename)
	if err != nil {
		return true
	}
	defer file.Close()

	buffer := make([]byte, 256)
	n, err := file.Read(buffer)
	if err != nil {
		return true
	}

	for i := 0; i < n; i++ {
		if buffer[i] == 0 {
			return true
		}
	}

	return false
}
