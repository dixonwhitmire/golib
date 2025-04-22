package iolib

import (
	"github.com/google/go-cmp/cmp"
	"os"
	"path"
	"testing"
)

const testDataDirectory = "./testdata"

// sampleTextFilePath returns the full path to the sample text file fixture.
func sampleTextFilePath(t *testing.T) string {
	t.Helper()
	return path.Join(testDataDirectory, "sample.txt")
}

// sampleCsvFilePath returns the full path toe the sample csv file fixture.
func sampleCsvFilePath(t *testing.T) string {
	t.Helper()
	return path.Join(testDataDirectory, "sample.csv")
}

// readSampleFile returns the content from our sample file as []byte.
func readSampleFile(t *testing.T) []byte {
	t.Helper()
	fileBytes, err := os.ReadFile(sampleTextFilePath(t))
	if err != nil {
		t.Fatalf("error reading sample file: %v", err)
	}
	return fileBytes
}

func TestReadFileContent_String(t *testing.T) {
	want := string(readSampleFile(t))

	got, err := ReadFileContent[string](sampleTextFilePath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("ReadFileContent[string]: %s", diff)
	}
}

func TestReadFileContent_Byte(t *testing.T) {
	want := readSampleFile(t)

	got, err := ReadFileContent[[]byte](sampleTextFilePath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("ReadFileContent[[]byte]]: %s", diff)
	}
}

func TestCsvRecord(t *testing.T) {
	want := make([][]string, 0)
	want = append(want, []string{"first_name", "last_name"})
	want = append(want, []string{"John", "Doe"})
	want = append(want, []string{"Jane", "Doe"})

	got := make([][]string, 0)

	seq, err := CsvRecord(sampleCsvFilePath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for csvRecord, err := range seq {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got = append(got, csvRecord)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("CsvRecord Diffs: %s", diff)
	}
}

func TestCsvRecord_ReadHeader(t *testing.T) {
	want := make([][]string, 0)
	want = append(want, []string{"first_name", "last_name"})

	got := make([][]string, 0)

	seq, err := CsvRecord(sampleCsvFilePath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for csvRecord, err := range seq {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got = append(got, csvRecord)
		break
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("CsvRecord Diffs: %s", diff)
	}
}

func TestFileLines(t *testing.T) {
	want := make([]string, 0)
	want = append(want, "This is a sample file with not much")
	want = append(want, "content at all!")

	got := make([]string, 0)

	seq, err := FileLines(sampleTextFilePath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for line := range seq {
		got = append(got, line)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Sample File Diffs: %s", diff)
	}
}

func TestFileLines_OneLine(t *testing.T) {
	want := make([]string, 0)
	want = append(want, "This is a sample file with not much")

	got := make([]string, 0)

	seq, err := FileLines(sampleTextFilePath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for line := range seq {
		got = append(got, line)
		break
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Sample File Diffs: %s", diff)
	}
}

func TestCsvWriter(t *testing.T) {
	// configure writer for testing
	testCsvPath := path.Join(t.TempDir(), "test.csv")

	writer, err := NewCsvWriter(testCsvPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer writer.Close()

	// csv contents
	csvData := [][]string{{"first_name", "last_name"}, {"John", "Doe"}}
	for _, row := range csvData {
		err := writer.Write(row)
		if err != nil {
			t.Fatalf("unexpected error writing record: %v", err)
		}
	}
	writer.Flush()
}
