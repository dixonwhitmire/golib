// Package csvlib provides supports reading from and writing to CSV files using iterators and structs.
package csvlib

import (
	"encoding/csv"
	"fmt"
	"github.com/dixonwhitmire/golib/errorlib"
	"io"
	"iter"
	"os"
)

// Record encapsulates a csv record including the record's line number and associated data.
type Record[T any] struct {
	LineNumber int
	Data       T
}

// ConversionFunc is used to convert a CSV record []string to type T
type ConversionFunc[T any] func([]string) (T, error)

// Iterator returns a Seq2 iterator, returning Record[T] and an error, if applicable.
// Data from the underlying CSV file is mapped to T using a ConversionFunc.
func Iterator[T any](
	inputFilePath string,
	hasHeader bool,
	conversionFunc ConversionFunc[T]) (iter.Seq2[Record[T], error], error) {

	f, err := os.Open(inputFilePath)
	if err != nil {
		return nil, errorlib.CreateError("Iterator", fmt.Sprintf("error accessing: %s %v", inputFilePath, err))
	}

	if conversionFunc == nil {
		return nil, errorlib.CreateError("Iterator", "conversionFunc is required")
	}

	reader := csv.NewReader(f)

	return func(yield func(Record[T], error) bool) {
		defer f.Close()

		lineNumber := 0
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			lineNumber++
			if hasHeader && lineNumber == 1 {
				continue
			}
			convertedData, err := conversionFunc(record)
			rec := Record[T]{LineNumber: lineNumber, Data: convertedData}
			if !yield(rec, err) {
				break
			}
		}
	}, nil
}
