// Package iolib provides iterators for reading text based file data.
package iolib

import (
	"bufio"
	"fmt"
	"github.com/dixonwhitmire/golib/errorlib"
	"iter"
	"os"
)

// FileContent is the type set used as a return type in ReadFileContents
type FileContent interface {
	~string | ~[]byte
}

// ReadFileContent returns the contents of a file as FileContent.
func ReadFileContent[T FileContent](filePath string) (T, error) {
	var contentType T

	fileBytes, err := os.ReadFile(filePath)

	if err != nil {
		return contentType, errorlib.CreateError("ReadFileContent", fmt.Sprintf("could not read file:%q, %v", filePath, err))
	}

	switch any(contentType).(type) {
	case string:
		return any(string(fileBytes)).(T), nil
	case []byte:
		return any(fileBytes).(T), nil
	default:
		return contentType, errorlib.CreateError("ReadFileContent", fmt.Sprintf("unknown content type:%T", any(fileBytes)))
	}
}

// FileLinesIterator returns a Seq containing a single file line.
func FileLinesIterator(csvFilePath string) (iter.Seq[string], error) {
	file, err := os.Open(csvFilePath)
	if err != nil {
		return nil, errorlib.CreateError("FileLinesIterator", fmt.Sprintf("could not open file:%q %v", csvFilePath, err))
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
