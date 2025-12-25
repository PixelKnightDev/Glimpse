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
	Content []string  // File lines (original)
	Lower   []string  // File lines (lowercase for fast case-insensitive search)
	ModTime time.Time // Last modification time
}

// Index represents an in-memory index of files in a directory
type Index struct {
	files map[string]*FileEntry
	mu    sync.RWMutex
}

// BuildIndex creates an index of all text files in the given directory
// progressCallback is called with current and total during indexing
func BuildIndex(dir string, progressCallback func(current, total int)) *Index {
	index := &Index{
		files: make(map[string]*FileEntry),
	}

	// First pass count total files for progress reporting
	totalFiles := countFiles(dir)

	// Second pass index all files
	currentCount := 0

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip directories that should be excluded
		if info.IsDir() {
			name := info.Name()
			if name == ".git" || name == "node_modules" || name == ".vscode" ||
				name == "target" || name == "build" || name == "dist" ||
				name == "venv" || name == ".venv" || name == "env" || name == ".env" ||
				name == "__pycache__" || name == ".pytest_cache" || name == "site-packages" {
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
				name == "target" || name == "build" || name == "dist" ||
				name == "venv" || name == ".venv" || name == "env" || name == ".env" ||
				name == "__pycache__" || name == ".pytest_cache" || name == "site-packages" {
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
	var lower []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		lower = append(lower, strings.ToLower(line))
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	idx.mu.Lock()
	defer idx.mu.Unlock()

	idx.files[path] = &FileEntry{
		Path:    path,
		Content: lines,
		Lower:   lower,
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
		for lineNum := 0; lineNum < len(entry.Content); lineNum++ {
			if fileResults >= 5 {
				break
			}

			// Use pre-lowercased content for case-insensitive search (was causing big lag without)
			searchLine := entry.Content[lineNum]
			if options.CaseInsensitive {
				searchLine = entry.Lower[lineNum]
			}

			if strings.Contains(searchLine, searchPattern) {
				results = append(results, Result{
					File:    entry.Path,
					Line:    lineNum + 1, //  because line numbers are 1-indexed
					Content: entry.Content[lineNum],
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

// returns the cached content lines for a file
// Returns nil if file not found in index
func (idx *Index) GetFileContent(filename string) []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	if entry, exists := idx.files[filename]; exists {
		return entry.Content
	}
	return nil
}
