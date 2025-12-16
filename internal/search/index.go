/*
 * Glimpse - Index Module
 * Author: Pratyush Yadav <pratyushyadav0106@gmail.com>
 * Description: In-memory indexing system for fast code search
 * License: MIT License
 */
package search

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// FileEntry represents a single indexed file with its content
type FileEntry struct {
	Path    string    // File path
	Content []string  // File lines
	ModTime time.Time // Last modification time
}

// Index represents an in-memory index of files in a directory
type Index struct {
	files map[string]*FileEntry
	mu    sync.RWMutex
}

// BuildIndex creates an index of all text files in the given directory
// progressCallback is called with (current, total) during indexing
func BuildIndex(dir string, progressCallback func(current, total int)) *Index {
	index := &Index{
		files: make(map[string]*FileEntry),
	}

	// First pass: count total files for progress reporting
	totalFiles := countFiles(dir)

	// Second pass: index all files
	currentCount := 0

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip directories that should be excluded
		if info.IsDir() {
			name := info.Name()
			if name == ".git" || name == "node_modules" || name == ".vscode" ||
				name == "target" || name == "build" || name == "dist" {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip binary files
		if IsBinaryFile(path) {
			return nil
		}

		// Add file to index
		index.addFile(path)
		currentCount++

		// Call progress callback
		if progressCallback != nil {
			progressCallback(currentCount, totalFiles)
		}

		return nil
	})

	return index
}

// countFiles counts total files in directory (for progress reporting)
func countFiles(dir string) int {
	count := 0

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip directories that should be excluded
		if info.IsDir() {
			name := info.Name()
			if name == ".git" || name == "node_modules" || name == ".vscode" ||
				name == "target" || name == "build" || name == "dist" {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip binary files
		if IsBinaryFile(path) {
			return nil
		}

		count++
		return nil
	})

	return count
}

// addFile reads a file and adds it to the index
func (idx *Index) addFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	idx.mu.Lock()
	defer idx.mu.Unlock()

	idx.files[path] = &FileEntry{
		Path:    path,
		Content: lines,
		ModTime: info.ModTime(),
	}

	return nil
}

// Search searches the index for a pattern and returns matching results
func (idx *Index) Search(pattern string, options SearchOptions) []Result {
	var results []Result

	maxResults := options.MaxResults
	if maxResults == 0 {
		maxResults = 50
	}

	searchPattern := pattern
	if options.CaseInsensitive {
		searchPattern = strings.ToLower(pattern)
	}

	idx.mu.RLock()
	defer idx.mu.RUnlock()

	// Iterate through indexed files
	for _, entry := range idx.files {
		if len(results) >= maxResults {
			break
		}

		// Search within this file (max 5 results per file)
		fileResults := 0
		for lineNum, line := range entry.Content {
			if fileResults >= 5 {
				break
			}

			searchLine := line
			if options.CaseInsensitive {
				searchLine = strings.ToLower(line)
			}

			if strings.Contains(searchLine, searchPattern) {
				results = append(results, Result{
					File:    entry.Path,
					Line:    lineNum + 1, // Line numbers are 1-indexed
					Content: line,
				})
				fileResults++
			}

			// Stop if we've reached max results
			if len(results) >= maxResults {
				break
			}
		}
	}

	return results
}

// FileCount returns the number of indexed files
func (idx *Index) FileCount() int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return len(idx.files)
}
