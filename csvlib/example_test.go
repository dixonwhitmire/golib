package csvlib

import (
	"fmt"
	"strings"
)

type ExampleRecord struct {
	FirstName string
	LastName  string
}

func conversionFunc(rec []string) (ExampleRecord, error) {
	return ExampleRecord{FirstName: rec[0], LastName: rec[1]}, nil
}

func ExampleNewDefaultIterator() {
	const csvData = `"first_name","last_name"
"John","Doe"
"Jane","Doe"
`
	input := strings.NewReader(csvData)
	iter, err := NewDefaultIterator[ExampleRecord](input, true, conversionFunc)
	if err != nil {
		fmt.Println("error: ", err.Error())
	}

	for rec, err := range iter {
		if err != nil {
			fmt.Println("error: ", err.Error())
			break
		}
		fmt.Printf("record: %d, first name: %q, last name: %q\n", rec.LineNumber, rec.Data.FirstName, rec.Data.LastName)
	}
	// Output:
	// record: 2, first name: "John", last name: "Doe"
	// record: 3, first name: "Jane", last name: "Doe"
}
