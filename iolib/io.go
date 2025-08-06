// Package iolib provides iterators for reading text based file data.
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
func FileLinesIterator(csvFilePath string) (iter.Seq[string], error) {
	file, err := os.Open(csvFilePath)
	if err != nil {
		return nil, fmt.Errorf("FileLinesIterator: could not open file:%q %w", csvFilePath, err)
	}

	scanner := bufio.NewScanner(file)
	return func(yield func(string) bool) {
		defer file.Close()
		for scanner.Scan() {
			line := scanner.Text()
			if !yield(line) {
				break
			}
		}
	}, nil
}
