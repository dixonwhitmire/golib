package iolib

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/dixonwhitmire/golib/errorlib"
	"io"
	"iter"
	"os"
)

// FileContent is the type set used as a return type in ReadFileContents
type FileContent interface {
	~string | ~[]byte
}

// ReadFileContent returns the contents of a file as FileContent.
func ReadFileContent[T FileContent](filePath string) (T, error) {
	// zero value for our type constraint
	var zero T

	fileBytes, err := os.ReadFile(filePath)

	if err != nil {
		return zero, errorlib.WrapError("ReadFileContent", fmt.Sprintf("could not read file:%s", filePath), err)
	}

	var contentType T

	switch any(contentType).(type) {
	case string:
		return any(string(fileBytes)).(T), nil
	case []byte:
		return any(fileBytes).(T), nil
	default:
		return zero, errorlib.WrapError("ReadFileContent", fmt.Sprintf("unknown content type:%T", any(fileBytes)), err)
	}
}

// CsvRecord encapsulates a csv record including the record's line number and associated data.
type CsvRecord struct {
	LineNumber int
	Data       []string
}

// CsvRecordIterator returns a Seq2 iterator which includes the CsvRecord and error (if applicable) encountered during the read operation.
func CsvRecordIterator(csvFilePath string) (iter.Seq2[CsvRecord, error], error) {
	csvFile, err := os.Open(csvFilePath)
	if err != nil {
		return nil, errorlib.WrapError("CsvRecordIterator", fmt.Sprintf("could not open file:%s", csvFilePath), err)
	}

	reader := csv.NewReader(csvFile)

	return func(yield func(CsvRecord, error) bool) {
		defer csvFile.Close()

		lineCounter := 0
		for {
			lineCounter++
			csvData, err := reader.Read()
			if errors.Is(err, io.EOF) {
				break
			}
			csvRecord := CsvRecord{lineCounter, csvData}
			if !yield(csvRecord, err) {
				break
			}
		}
	}, nil
}

// FileLinesIterator returns a Seq containing a single file line.
func FileLinesIterator(csvFilePath string) (iter.Seq[string], error) {
	file, err := os.Open(csvFilePath)
	if err != nil {
		return nil, errorlib.WrapError("FileLinesIterator", fmt.Sprintf("could not open file:%s", csvFilePath), err)
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

// Error returns the current errorlib
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
func NewCsvWriter(csvFilePath string) (CsvWriter, error) {
	f, err := os.Create(csvFilePath)
	if err != nil {
		return CsvWriter{}, errorlib.WrapError("NewCsvWriter", fmt.Sprintf("could not open file:%s", csvFilePath), err)
	}

	writer := CsvWriter{csvFile: f, csvWriter: csv.NewWriter(f)}
	return writer, nil
}

// parseCsvMetaData returns the header record (if present) for a csv file and a column count.
func parseCsvMetaData(csvFilePath string, hasHeader bool) ([]string, int, error) {
	var headerRecord []string
	columnCount := 0

	csvRecords, err := CsvRecordIterator(csvFilePath)
	if err != nil {
		return headerRecord, columnCount, errorlib.WrapError("parseCsvMetadata", "error creating CsvRecordIterator iterator", err)
	}

	for csvRecord, err := range csvRecords {
		if err != nil {
			return headerRecord, columnCount, errorlib.WrapError("parseCsvMetadata", "error parsing csv metadata", err)
		}
		if hasHeader {
			headerRecord = csvRecord.Data
		}
		columnCount = len(csvRecord.Data)
		break
	}
	if columnCount == 0 {
		errorText := fmt.Sprintf("parseCsvMetaData: invalid csv file:%s column count == 0", csvFilePath)
		return headerRecord, columnCount, errorlib.CreateError("parseCsvMetaData", errorText)
	}
	return headerRecord, columnCount, nil
}

// MergeCsvFiles merges the provided inputFiles into a single file located at outputCsvPath.
// The first file in inputFiles establishes the merge file's structure and header row if hasHeader is true
func MergeCsvFiles(outputCsvPath string, hasHeader bool, inputCsvFiles ...string) error {
	if inputCsvFiles == nil || len(inputCsvFiles) == 0 {
		return errorlib.CreateError("MergeCsvFiles", "inputCsvFiles is required")
	}

	// grab the first file for the header (if any) and column count
	headerRecord, columnCount, err := parseCsvMetaData(inputCsvFiles[0], hasHeader)
	if err != nil {
		return errorlib.WrapError("MergeCsvFiles", "error parsing metadata", err)
	}

	// prepare the output file
	writer, err := NewCsvWriter(outputCsvPath)
	if err != nil {
		return errorlib.WrapError("MergeCsvFiles", "error creating CsvWriter", err)
	}
	// cleanup
	defer func() {
		writer.Flush()
		writer.Close()
	}()

	if hasHeader {
		err := writer.Write(headerRecord)
		if err != nil {
			return errorlib.WrapError("MergeCsvFiles", "error writing header", err)
		}
		writer.Flush()
	}

	// process each file
	for _, csvFilePath := range inputCsvFiles {
		lineCounter := 0

		// create our iterator
		csvRecords, err := CsvRecordIterator(csvFilePath)
		if err != nil {
			return errorlib.WrapError("MergeCsvFiles", "error creating CsvRecordIterator iterator", err)
		}

		for csvRecord, err := range csvRecords {
			lineCounter++
			if err != nil {
				return errorlib.WrapError("MergeCsvFiles", "error reading csv record", err)
			}

			// first line is used to compare column counts and if a header row, skip
			if lineCounter == 1 {
				if len(csvRecord.Data) != columnCount {
					errorText := fmt.Sprintf("invalid csv column count. expected=%d, actual=%d file=%s",
						columnCount, len(csvRecord.Data), csvFilePath)
					return errorlib.CreateError("MergeCsvFiles", errorText)
				}
				if hasHeader {
					continue
				}
			}
			err := writer.Write(csvRecord.Data)
			if err != nil {
				return errorlib.WrapError(
					"MergeCsvFiles", fmt.Sprintf("error writing csv record line:%d", lineCounter), err)
			}
			writer.Flush()
		}
	}
	return nil
}
