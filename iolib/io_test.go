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

// readSampleFile returns the content from our sample file as []byte.
func readSampleFile(t *testing.T) []byte {
	t.Helper()
	fileBytes, err := os.ReadFile(sampleTextFilePath(t))
	if err != nil {
		t.Fatalf("readSampleFile unexpected error: %v", err)
	}
	return fileBytes
}

func TestReadFileContent_String(t *testing.T) {
	want := string(readSampleFile(t))

	got, err := ReadFileContent[string](sampleTextFilePath(t))
	if err != nil {
		t.Fatalf("ReadFileContent unexpected error: %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("ReadFileContent found diff (-want +got):\n%s", diff)
	}
}

func TestReadFileContent_Byte(t *testing.T) {
	want := readSampleFile(t)

	got, err := ReadFileContent[[]byte](sampleTextFilePath(t))
	if err != nil {
		t.Fatalf("ReadFileContent unexpected error: %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("ReadFileContent found diff (-want +got):\n%s", diff)
	}
}

func TestFileLinesIterator(t *testing.T) {
	want := make([]string, 0)
	want = append(want, "This is a sample file with not much")
	want = append(want, "content at all!")

	got := make([]string, 0)

	seq, err := FileLinesIterator(sampleTextFilePath(t))
	if err != nil {
		t.Fatalf("FileLinesIterator unexpected error: %v", err)
	}

	for line := range seq {
		got = append(got, line)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("FileLinesIterator found diff (-want +got):\n%s", diff)
	}
}

func TestFileLinesIterator_OneLine(t *testing.T) {
	want := make([]string, 0)
	want = append(want, "This is a sample file with not much")

	got := make([]string, 0)

	seq, err := FileLinesIterator(sampleTextFilePath(t))
	if err != nil {
		t.Fatalf("FileLinesIterator unexpected error: %v", err)
	}

	for line := range seq {
		got = append(got, line)
		break
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("FileLinesIterator found diff (-want +got):\n%s", diff)
	}
}
