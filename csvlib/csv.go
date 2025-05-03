// Package csvlib provides supports reading from and writing to CSV files using iterators and structs.
package csvlib

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/dixonwhitmire/golib/errorlib"
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
	// cause is the underlying exception if applicable.
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

// ConversionError is returned when an error occurs converting the underlying CSV record from []string to type T.
type ConversionError struct {
	*baseError
}

// Error returns context specific ConversionError information.
func (ce *ConversionError) Error() string {
	return ce.format("ConversionError")
}

// Unwrap returns our inner error, aka "the cause".
func (ce *ConversionError) Unwrap() error {
	return ce.cause
}

// NewConversionError returns a ConversionError with the specified context information.
func NewConversionError(filePath string, lineNumber int, cause error) *ConversionError {
	return &ConversionError{baseError: &baseError{filePath: filePath, lineNumber: lineNumber, cause: cause}}
}

// Record encapsulates a csv record including the record's line number and associated data.
type Record[T any] struct {
	LineNumber int
	Data       T
}

// ConversionFunc is used to convert a CSV record []string to type T
type ConversionFunc[T any] func([]string) (T, error)

// NewDefaultIterator returns an iterator with a DefaultBufferSize.
func NewDefaultIterator[T any](inputFilePath string,
	hasHeader bool,
	conversionFunc ConversionFunc[T]) (iter.Seq2[Record[T], error], error) {
	return iterator[T](inputFilePath, hasHeader, DefaultBufferSize, conversionFunc)
}

// NewIterator returns a buffered iterator with a configurable buffer size.
// DefaultBufferSize is used if bufferSize < DefaultBufferSize.
func NewIterator[T any](inputFilePath string,
	hasHeader bool,
	bufferSize int,
	conversionFunc ConversionFunc[T]) (iter.Seq2[Record[T], error], error) {
	return iterator[T](inputFilePath, hasHeader, bufferSize, conversionFunc)
}

// iterator returns iter.Seq2[Record[T], error.
// Data from the underlying CSV file is read using a buffered csv.Reader and is mapped to T using a ConversionFunc.
// Custom buffer sizes may be specified if customBufferSize is set to a value > DefaultBufferSize.
func iterator[T any](
	inputFilePath string,
	hasHeader bool,
	customBufferSize int,
	conversionFunc ConversionFunc[T]) (iter.Seq2[Record[T], error], error) {

	if conversionFunc == nil {
		err := errorlib.CreateError("iterator", "conversionFunc is required")
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
				if !yield(Record[T]{LineNumber: lineNumber}, NewConversionError(inputFilePath, lineNumber, err)) {
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
