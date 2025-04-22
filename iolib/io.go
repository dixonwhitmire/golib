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

// CsvWriter writes CSV records to an output file.
// Use NewCsvWriter to create new instances and Close to close the underlying file.
type CsvWriter struct {
	csvFile   *os.File
	csvWriter *csv.Writer
}

// Error returns the current error
func (w *CsvWriter) Error() error {
	return w.csvWriter.Error()
}

// Flush writes unbuffered data to disk.
func (w *CsvWriter) Flush() {
	w.csvWriter.Flush()
}

// Write writes a csvRecord to the underlying output file.
func (w *CsvWriter) Write(csvRecord []string) error {
	return w.csvWriter.Write(csvRecord)
}

// Close closes the underlying CsvWriter file.
func (w *CsvWriter) Close() error {
	return w.csvFile.Close()
}

// NewCsvWriter creates a new CsvWriter instance.
func NewCsvWriter(filePath string) (CsvWriter, error) {
	f, err := os.Create(filePath)
	if err != nil {
		return CsvWriter{}, fmt.Errorf("could not open file: %w", err)
	}

	writer := CsvWriter{csvFile: f, csvWriter: csv.NewWriter(f)}
	return writer, nil
}
