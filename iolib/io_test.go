package iolib

import (
	"github.com/google/go-cmp/cmp"
	"os"
	"testing"
)

const sampleTextFilePath = "./testdata/sample.txt"

// readSampleFile returns the content from our sample file as []byte.
func readSampleFile(t *testing.T) []byte {
	t.Helper()
	fileBytes, err := os.ReadFile(sampleTextFilePath)
	if err != nil {
		t.Fatalf("readSampleFile unexpected error: %v", err)
	}
	return fileBytes
}

func TestReadFileAsString(t *testing.T) {
	want := string(readSampleFile(t))

	got, err := ReadFileAsString(sampleTextFilePath)
	if err != nil {
		t.Fatalf("ReadFileAsString unexpected error: %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("ReadFileAsString found diff (-want +got):\n%s", diff)
	}
}

func TestReadFileAsBytes(t *testing.T) {
	want := readSampleFile(t)

	got, err := ReadFileAsBytes(sampleTextFilePath)
	if err != nil {
		t.Fatalf("ReadFileAsBytes unexpected error: %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("ReadFileAsBytes found diff (-want +got):\n%s", diff)
	}
}

func TestFileLinesIterator(t *testing.T) {
	want := make([]string, 0)
	want = append(want, "This is a sample file with not much")
	want = append(want, "content at all!")

	got := make([]string, 0)

	seq, err := FileLinesIterator(sampleTextFilePath)
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

	seq, err := FileLinesIterator(sampleTextFilePath)
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
