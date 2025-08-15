// Package iolib provides access to text file contents in memory or via iterators.
package iolib

import (
	"bufio"
	"fmt"
	"iter"
	"os"
)

// ReadFileAsString returns the contents of a file as a string.
func ReadFileAsString(filePath string) (string, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("ReadFileAsString: could not read file %q: %w", filePath, err)
	}
	return string(fileBytes), nil
}

// ReadFileAsBytes returns the contents of a file as a byte slice.
func ReadFileAsBytes(filePath string) ([]byte, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ReadFileAsBytes: could not read file %q: %w", filePath, err)
	}
	return fileBytes, nil
}

// FileLinesIterator returns a Seq containing a single file line.
func FileLinesIterator(filePath string) (iter.Seq2[string, error], error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("FileLinesIterator: could not open file:%q %w", filePath, err)
	}

	// iterator function includes an error if applicable
	return func(yield func(string, error) bool) {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if !yield(line, nil) {
				break
			}
		}

		// follow-up to check for errors
		if scanErr := scanner.Err(); scanErr != nil {
			yield("", fmt.Errorf("FileLinesIterator: scan error: %w", scanErr))
		}

		if closeErr := file.Close(); closeErr != nil {
			yield("", fmt.Errorf("FileLinesIterator: close error: %w", closeErr))
		}
	}, nil
}
