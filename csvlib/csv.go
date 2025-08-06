// Package csvlib provides supports reading from and writing to CSV files using iterators and structs.
package csvlib

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"iter"
	"os"
)

// DefaultBufferSize is the default buffer size used for reading records.
const DefaultBufferSize = 1024 * 4

// baseError provides the general structure for the package's exported errors.
// The filePath and lineNumber, if applicable, provide the general context for the error.
type baseError struct {
	// filePath is the file system path to the csv file.
	filePath string
	// lineNumber is the line number where the error occurred.
	// lineNumber == 0 for errors not related to error processing.
	lineNumber int
	// operation specifies the process iteration, record parsing, etc. which received an error.
	operation string
	// cause is the underlying error if available.
	cause error
}

// format returns the error string for the specified operation.
// The returned error string includes the operation and file path, and optionally includes the line number and error, if
// present.
func (b *baseError) format(operation string) string {
	var message string

	if b.lineNumber > 0 {
		message = fmt.Sprintf("%s error in %s line %d", operation, b.filePath, b.lineNumber)
	} else {
		message = fmt.Sprintf("%s error in %s", operation, b.filePath)
	}

	if b.cause != nil {
		message += fmt.Sprintf(" %v", b.cause)
	}
	return message
}

// IterationError is returned when an error occurs during CSV file iteration.
type IterationError struct {
	*baseError
}

// Error returns context specific IterationError information.
func (ie *IterationError) Error() string {
	return ie.format("IterationError")
}

// Unwrap returns our inner error, aka "the cause".
func (ie *IterationError) Unwrap() error {
	return ie.cause
}

func NewIterationError(filePath string, lineNumber int, cause error) *IterationError {
	return &IterationError{baseError: &baseError{filePath: filePath, lineNumber: lineNumber, cause: cause}}
}

// ParseError is returned when an error occurs parsing a raw []string record to type T.
type ParseError struct {
	*baseError
}

// Error returns context specific ParseError information.
func (ce *ParseError) Error() string {
	return ce.format("ParseError")
}

// Unwrap returns our inner error, aka "the cause".
func (ce *ParseError) Unwrap() error {
	return ce.cause
}

// NewParseError returns a ParseError with the specified context information.
func NewParseError(filePath string, lineNumber int, cause error) *ParseError {
	return &ParseError{baseError: &baseError{filePath: filePath, lineNumber: lineNumber, cause: cause}}
}

// Record encapsulates a csv record including the record's line number and associated data.
type Record[T any] struct {
	LineNumber int
	Data       T
}

// ParseFunc is used to parse a CSV record []string into type T
type ParseFunc[T any] func([]string) (T, error)

// NewDefaultIterator returns an iterator with a DefaultBufferSize.
func NewDefaultIterator[T any](inputFilePath string,
	hasHeader bool,
	conversionFunc ParseFunc[T]) (iter.Seq2[Record[T], error], error) {
	return iterator[T](inputFilePath, hasHeader, DefaultBufferSize, conversionFunc)
}

// NewIterator returns a buffered iterator with a configurable buffer size.
// DefaultBufferSize is used if bufferSize < DefaultBufferSize.
func NewIterator[T any](inputFilePath string,
	hasHeader bool,
	bufferSize int,
	conversionFunc ParseFunc[T]) (iter.Seq2[Record[T], error], error) {
	return iterator[T](inputFilePath, hasHeader, bufferSize, conversionFunc)
}

// iterator returns iter.Seq2[Record[T], error].
// Data from the underlying CSV file is read using a buffered csv.Reader and is mapped to T using a ParseFunc.
// Custom buffer sizes may be specified if customBufferSize is set to a value > DefaultBufferSize.
func iterator[T any](
	inputFilePath string,
	hasHeader bool,
	customBufferSize int,
	conversionFunc ParseFunc[T]) (iter.Seq2[Record[T], error], error) {

	if conversionFunc == nil {
		err := errors.New("iterator: conversionFunc is required")
		return nil, NewIterationError(inputFilePath, 0, err)
	}

	f, err := os.Open(inputFilePath)
	if err != nil {
		return nil, NewIterationError(inputFilePath, 0, err)
	}

	bufSize := DefaultBufferSize
	if customBufferSize > DefaultBufferSize {
		bufSize = customBufferSize
	}
	reader := csv.NewReader(bufio.NewReaderSize(f, bufSize))

	return func(yield func(Record[T], error) bool) {
		defer func() { _ = f.Close() }()

		lineNumber := 0
		if hasHeader {
			lineNumber++
			if _, err := reader.Read(); err != nil {
				return
			}
		}

		for {
			lineNumber++
			csvFields, err := reader.Read()
			// handle csv read errors
			if err != nil {
				if err == io.EOF {
					return
				}
				// general iteration error
				if !yield(Record[T]{LineNumber: lineNumber}, NewIterationError(inputFilePath, lineNumber, err)) {
					return
				}
				continue
			}
			convertedData, err := conversionFunc(csvFields)
			// record conversion error
			if err != nil {
				if !yield(Record[T]{LineNumber: lineNumber}, NewParseError(inputFilePath, lineNumber, err)) {
					return
				}
				continue
			}
			if !yield(Record[T]{LineNumber: lineNumber, Data: convertedData}, nil) {
				return
			}
		}
	}, nil
}

// ConvertFunc converts type T to a []string record for CSV writing output.
type ConvertFunc[T any] func(T) ([]string, error)

// Writer writes csv records to an output file.
// convertFunc is used to convert type T to []string for CSV writer output.
type Writer[T any] struct {
	convertFunc  ConvertFunc[T]
	outputFile   *os.File
	outputWriter *csv.Writer
}

// Flush writes the current buffer to the output
func (w *Writer[T]) Flush() {
	w.outputWriter.Flush()
}

// Close flushes remaining data to the writer and closes related resources.
func (w *Writer[T]) Close() error {
	w.outputWriter.Flush()
	err := w.outputWriter.Error()
	if err != nil {
		return fmt.Errorf("Writer.Close: error flushing data %w", err)
	}
	err = w.outputFile.Close()
	if err != nil {
		return fmt.Errorf("Writer.Close: error closing file %w", err)
	}
	return nil
}

// Write writes the csv record to the underlying output file.
func (w *Writer[T]) Write(inputRecord T) error {
	csvFields, err := w.convertFunc(inputRecord)
	if err != nil {
		return fmt.Errorf("Writer.Write: error converting %w", err)
	}

	err = w.outputWriter.Write(csvFields)
	if err != nil {
		return fmt.Errorf("Writer.Write: error writing data %w", err)
	}
	return nil
}

// WriteHeader writes a header to the output csv file.
func (w *Writer[T]) WriteHeader(headerRecord []string) error {
	err := w.outputWriter.Write(headerRecord)
	if err != nil {
		return fmt.Errorf("Writer.WriteHeade: error writing header %w", err)
	}
	return nil
}

// NewDefaultWriter creates a new Writer which writes records of type T to an output CSV file.
func NewDefaultWriter[T any](outputFilePath string, convertFunc ConvertFunc[T]) (Writer[T], error) {
	if convertFunc == nil {
		return Writer[T]{}, errors.New("NewDefaultWriter: conversionFunc is required")
	}

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return Writer[T]{}, fmt.Errorf("NewDefaultWriter: error creating output file %w", err)
	}
	writer := csv.NewWriter(bufio.NewWriterSize(outputFile, DefaultBufferSize))
	return Writer[T]{convertFunc: convertFunc, outputFile: outputFile, outputWriter: writer}, nil
}
