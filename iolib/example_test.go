package iolib

import "fmt"

const filePath = "testdata/sample.txt"

func ExampleReadFileAsString() {
	fileContents, err := ReadFileAsString(filePath)
	if err != nil {
		fmt.Println("error reading file ", err.Error())
		return
	}

	fmt.Println(fileContents)
	// Output:
	// This is a sample file with not much
	// content at all!
}

func ExampleReadFileAsBytes() {
	fileContents, err := ReadFileAsBytes(filePath)
	if err != nil {
		fmt.Println("error reading file ", err.Error())
		return
	}

	fmt.Println(fileContents)
	// Output:
	// [84 104 105 115 32 105 115 32 97 32 115 97 109 112 108 101 32 102 105 108 101 32 119 105 116 104 32 110 111 116 32 109 117 99 104 10 99 111 110 116 101 110 116 32 97 116 32 97 108 108 33]
}
