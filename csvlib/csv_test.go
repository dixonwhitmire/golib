package csvlib

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"strings"
	"testing"
)

const (
	sampleFilePath         = "./testdata/sample.csv"
	sampleFileNoHeaderPath = "./testdata/sample-no-header.csv"
	sampleFileErrorPath    = "./testdata/sample-conversion-error.csv"
)

// CustomRecord is used in test cases as a "concrete" value for our generic type parameter.
// CustomRecord is an "exported type" so that we can use the cmp.Diff tooling without additional options.
type CustomRecord struct {
	FirstName string
	LastName  string
}

// customRecordConversionFunc is used to convert csv fields to a CustomRecord for test cases.
func customRecordConversionFunc(csvFields []string) (CustomRecord, error) {
	customRecord := CustomRecord{
		FirstName: csvFields[0],
		LastName:  csvFields[1],
	}
	if strings.EqualFold(customRecord.LastName, "error") {
		return customRecord, NewConversionError(sampleFileErrorPath, 2, errors.New("test case error"))
	}
	return customRecord, nil
}

// iteratorTestCase uses a type parameter to parameterize iterator test cases.
type iteratorTestCase[T any] map[string]struct {
	csvPath   string
	hasHeader bool
	want      []Record[T]
	wantErr   bool
}

// runIteratorTestCases executes the iteratorTestCase for the specified parameter type
func runIteratorTestCases[T any](t *testing.T, conv ConversionFunc[T], cases iteratorTestCase[T]) {
	t.Helper()

	for name, tt := range cases {
		t.Run(name, func(subT *testing.T) {
			got := make([]Record[T], 0)
			iter, err := NewDefaultIterator[T](tt.csvPath, tt.hasHeader, conv)
			if (err != nil) != tt.wantErr {
				t.Fatalf("iterator error = %v, wantErr %v", err, tt.wantErr)
			}
			for rec, err := range iter {
				if err != nil {
					t.Errorf("iterator error = %v", err)
				}
				got = append(got, rec)
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("iterator found diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestIterator_StringSlice(t *testing.T) {
	tt := iteratorTestCase[[]string]{
		"has-header": {
			csvPath:   sampleFilePath,
			hasHeader: true,
			want: []Record[[]string]{
				{LineNumber: 2, Data: []string{"John", "Doe"}},
				{LineNumber: 3, Data: []string{"Jane", "Doe"}},
			},
			wantErr: false,
		},
		"has-no-header": {
			csvPath:   sampleFileNoHeaderPath,
			hasHeader: false,
			want: []Record[[]string]{
				{LineNumber: 1, Data: []string{"Steve", "Doe"}},
				{LineNumber: 2, Data: []string{"Sally", "Doe"}},
			},
			wantErr: false,
		},
	}

	conversionFunc := func(records []string) ([]string, error) {
		return records, nil
	}

	runIteratorTestCases[[]string](t, conversionFunc, tt)
}

func TestIterator_CustomType(t *testing.T) {
	tt := iteratorTestCase[CustomRecord]{
		"has-header": {
			csvPath:   sampleFilePath,
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
			csvPath:   sampleFileNoHeaderPath,
			hasHeader: false,
			want: []Record[CustomRecord]{
				{
					LineNumber: 1,
					Data: CustomRecord{
						FirstName: "Steve",
						LastName:  "Doe",
					},
				},
				{
					LineNumber: 2,
					Data: CustomRecord{
						FirstName: "Sally",
						LastName:  "Doe",
					},
				},
			},
			wantErr: false,
		},
	}
	runIteratorTestCases[CustomRecord](t, customRecordConversionFunc, tt)
}

func TestIterator_IterationError(t *testing.T) {
	_, err := NewDefaultIterator[CustomRecord]("/tmp/not-a-real.csv", true, customRecordConversionFunc)
	if err == nil {
		t.Fatal("NewDefaultIterator did not return an IterationError")
	} else {
		var ie *IterationError
		if !errors.As(err, &ie) {
			t.Errorf("NewDefaultIterator did not return an IterationError got %v", err)
		}
	}
}

func TestIterator_ConversionError(t *testing.T) {
	iter, err := NewDefaultIterator[CustomRecord](sampleFileErrorPath, true, customRecordConversionFunc)
	if err != nil {
		t.Fatalf("NewDefaultIterator unexpected error %v", err)
	}

	var ce *ConversionError
	for _, err := range iter {
		if !errors.As(err, &ce) {
			t.Errorf("NewDefaultIterator did not return an ConversionError got %v", err)
		}
	}
}
