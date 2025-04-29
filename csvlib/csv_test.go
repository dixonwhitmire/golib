package csvlib

import (
	"github.com/google/go-cmp/cmp"
	"path"
	"testing"
)

const testDataDirectory = "./testdata"

func TestIterator_StringSlice(t *testing.T) {
	tt := map[string]struct {
		inputFilePath string
		hasHeader     bool
		expected      []Record[[]string]
	}{
		"has-header": {
			inputFilePath: path.Join(testDataDirectory, "sample.csv"),
			hasHeader:     true,
			expected: []Record[[]string]{
				{LineNumber: 2, Data: []string{"John", "Doe"}},
				{LineNumber: 3, Data: []string{"Jane", "Doe"}},
			},
		},
		"no-header": {
			inputFilePath: path.Join(testDataDirectory, "sample-no-header.csv"),
			hasHeader:     false,
			expected: []Record[[]string]{
				{LineNumber: 1, Data: []string{"Steve", "Doe"}},
				{LineNumber: 2, Data: []string{"Sally", "Doe"}},
			},
		},
	}

	// conversionFunc is our pass-through conversion function which returns the csv record []string
	conversionFunc := func(record []string) ([]string, error) {
		return record, nil
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {

			got := make([]Record[[]string], 0)
			recs, err := Iterator[[]string](tc.inputFilePath, tc.hasHeader, conversionFunc)

			if err != nil {
				t.Fatalf("Iterator(%q, %t) failed: %v", tc.inputFilePath, tc.hasHeader, err)
			}

			for rec, err := range recs {
				if err != nil {
					t.Fatalf("Iterator(%q, %t) failed: %v", tc.inputFilePath, tc.hasHeader, err)
				}
				got = append(got, rec)
			}

			if diff := cmp.Diff(tc.expected, got); diff != "" {
				t.Errorf("Iterator(%q, %t) found diff (-want +got):\n%s", tc.inputFilePath, tc.hasHeader, diff)
			}
		})
	}
}
