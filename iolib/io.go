package iolib

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"iter"
	"os"
)

// FileContent is the type set used in ReadFileContents
type FileContent interface {
	~string | ~[]byte
}

// ReadFileContent returns the contents of a file as FileContent.
func ReadFileContent[T FileContent](filePath string) (T, error) {
	// zero is the zero value for our type constraint
	var zero T

	fileBytes, err := os.ReadFile(filePath)

	if err != nil {
		return zero, fmt.Errorf("could not read file: %w", err)
	}

	var contentType T

	switch any(contentType).(type) {
	case string:
		return any(string(fileBytes)).(T), nil
	case []byte:
		return any(fileBytes).(T), nil
	default:
		return zero, fmt.Errorf("unknown content type: %T", any(fileBytes))
	}
}

// CsvRecord returns a Seq2 iterator which includes the csvRecord and any error encountered during the read operation.
func CsvRecord(filePath string) (iter.Seq2[[]string, error], error) {
	csvFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}

	reader := csv.NewReader(csvFile)

	return func(yield func([]string, error) bool) {
		defer csvFile.Close()

		for {
			csvRecord, err := reader.Read()
			if errors.Is(err, io.EOF) {
				break
			}
			if !yield(csvRecord, err) {
				break
			}
		}
	}, nil
}

// FileLines returns a Seq containing a single file line.
func FileLines(filePath string) (iter.Seq[string], error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
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
