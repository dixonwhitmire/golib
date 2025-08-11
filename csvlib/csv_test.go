package csvlib

import (
	"errors"
	"github.com/dixonwhitmire/golib/iolib"
	"github.com/google/go-cmp/cmp"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// sampleCsv is a helper function which returns the sample CSV payload used in unit tests.
func sampleCsv(t *testing.T, includeHeader bool) string {
	t.Helper()
	var csvData string
	if includeHeader {
		csvData = `"first_name","last_name"` + "\n"
	}
	csvData += `"John","Doe"` + "\n"
	csvData += `"Jane","Doe"`
	return csvData
}

// CustomRecord is used in test cases as a "concrete" value for our generic type parameter.
// CustomRecord is an "exported type" so that we can use the cmp.Diff tooling without additional options.
type CustomRecord struct {
	FirstName string
	LastName  string
}

// customRecordParseFunc is used to parse csv fields to a CustomRecord for csv reading/iterator test cases.
func customRecordParseFunc(csvFields []string) (CustomRecord, error) {
	customRecord := CustomRecord{
		FirstName: csvFields[0],
		LastName:  csvFields[1],
	}
	if strings.EqualFold(customRecord.LastName, "error") {
		return customRecord, NewParseError(2, errors.New("test case error"))
	}
	return customRecord, nil
}

// customRecordConvertFunc is used to convert a type T to []string for csv writing.
func customRecordConvertFunc(customRecord CustomRecord) ([]string, error) {
	return []string{customRecord.FirstName, customRecord.LastName}, nil
}

// iteratorTestCase uses a type parameter to parameterize iterator test cases.
type iteratorTestCase[T any] map[string]struct {
	hasHeader bool
	want      []Record[T]
	wantErr   bool
}

// runIteratorTestCases executes the iteratorTestCase for the specified parameter type
func runIteratorTestCases[T any](t *testing.T, conv ParseFunc[T], cases iteratorTestCase[T]) {
	t.Helper()

	for name, tt := range cases {
		t.Run(name, func(subT *testing.T) {
			input := strings.NewReader(sampleCsv(t, tt.hasHeader))

			got := make([]Record[T], 0, 2)

			iter, err := NewDefaultIterator[T](input, tt.hasHeader, conv)
			if (err != nil) != tt.wantErr {
				t.Fatalf("iterator error = %v, wantErr %v", err, tt.wantErr)
			}
			for rec, err := range iter {
				if err != nil {
					subT.Errorf("iterator error = %v", err)
				}
				got = append(got, rec)
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				subT.Errorf("iterator found diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestIterator_StringSlice(t *testing.T) {
	tt := iteratorTestCase[[]string]{
		"has-header": {
			hasHeader: true,
			want: []Record[[]string]{
				{LineNumber: 2, Data: []string{"John", "Doe"}},
				{LineNumber: 3, Data: []string{"Jane", "Doe"}},
			},
			wantErr: false,
		},
		"has-no-header": {
			hasHeader: false,
			want: []Record[[]string]{
				{LineNumber: 1, Data: []string{"John", "Doe"}},
				{LineNumber: 2, Data: []string{"Jane", "Doe"}},
			},
			wantErr: false,
		},
	}

	parseFunc := func(records []string) ([]string, error) {
		return records, nil
	}

	runIteratorTestCases[[]string](t, parseFunc, tt)
}

func TestIterator_CustomType(t *testing.T) {
	tt := iteratorTestCase[CustomRecord]{
		"has-header": {
			hasHeader: true,
			want: []Record[CustomRecord]{
				{
					LineNumber: 2,
					Data: CustomRecord{
						FirstName: "John",
						LastName:  "Doe",
					},
				},
				{
					LineNumber: 3,
					Data: CustomRecord{
						FirstName: "Jane",
						LastName:  "Doe",
					},
				},
			},
			wantErr: false,
		},
		"has-no-header": {
			hasHeader: false,
			want: []Record[CustomRecord]{
				{
					LineNumber: 1,
					Data: CustomRecord{
						FirstName: "John",
						LastName:  "Doe",
					},
				},
				{
					LineNumber: 2,
					Data: CustomRecord{
						FirstName: "Jane",
						LastName:  "Doe",
					},
				},
			},
			wantErr: false,
		},
	}
	runIteratorTestCases[CustomRecord](t, customRecordParseFunc, tt)
}

func TestIterator_ConversionError(t *testing.T) {
	// Our test conversion func raises an error if the string "error" is in the "last name" field.
	sampleCsvData := strings.Replace(sampleCsv(t, true), "Doe", "error", 1)
	input := strings.NewReader(sampleCsvData)

	iter, err := NewDefaultIterator[CustomRecord](input, true, customRecordParseFunc)
	if err != nil {
		t.Fatalf("NewDefaultIterator unexpected error %v", err)
	}

	var pe *ParseError
	foundErr := false
	for _, err := range iter {
		if err != nil {
			if errors.As(err, &pe) {
				foundErr = true
			} else {
				t.Errorf("NewDefaultIterator did not return a ParseError got %v", err)
			}
		}
	}

	if !foundErr {
		t.Fatal("Expected a ParseError but did not receive one")
	}
}

func TestWriter(t *testing.T) {
	expectedContents := "first_name,last_name\nJohn,Doe\n"
	outputFilePath := filepath.Join(t.TempDir(), "test-writer.csv")

	outputFile, err := os.Create(outputFilePath)
	var convertFunc = ConvertFunc[CustomRecord](customRecordConvertFunc)

	w, err := NewDefaultWriter[CustomRecord](outputFile, convertFunc)
	if err != nil {
		t.Fatalf("NewDefaultWriter unexpected error %v", err)
	}
	defer w.Close()

	err = w.WriteHeader([]string{"first_name", "last_name"})
	if err != nil {
		t.Fatalf("WriteHeader unexpected error %v", err)
	}
	err = w.Write(CustomRecord{
		FirstName: "John",
		LastName:  "Doe",
	})
	if err != nil {
		t.Fatalf("Write unexpected error %v", err)
	}
	w.Flush()

	actualOutput, err := iolib.ReadFileAsString(outputFilePath)
	if err != nil {
		t.Fatalf("ReadFileContent unexpected error %v", err)
	}

	if diff := cmp.Diff(expectedContents, actualOutput); diff != "" {
		t.Errorf("Writer did not write expected contents (-want +got):\n%s", diff)
	}
}
