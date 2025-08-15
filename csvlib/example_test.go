package csvlib

import (
	"bytes"
	"fmt"
	"strings"
)

// ExampleRecord is the parse target for the csv iterator.
type ExampleRecord struct {
	FirstName string
	LastName  string
}

// csvToExampleRecord concerts a csv record, []string, to our ExampleRecord.
func csvToExampleRecord(rec []string) (ExampleRecord, error) {
	return ExampleRecord{FirstName: rec[0], LastName: rec[1]}, nil
}

func ExampleNewDefaultIterator() {
	// prepare sample data
	const csvData = `"first_name","last_name"
"John","Doe"
"Jane","Doe"
`
	input := strings.NewReader(csvData)

	iter, err := NewDefaultIterator[ExampleRecord](input, true, csvToExampleRecord)
	if err != nil {
		fmt.Println("error: ", err.Error())
	}

	// process each record
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

// exampleRecordToCsv converts ExampleRecord to a []string/csv record.
func exampleRecordToCsv(rec ExampleRecord) ([]string, error) {
	return []string{rec.FirstName, rec.LastName}, nil
}

func ExampleNewWriter() {
	var output bytes.Buffer

	writer, err := NewWriter(&output, exampleRecordToCsv)
	if err != nil {
		fmt.Println("error creating writer ", err.Error())
		return
	}
	// writer will flush prior to close
	defer writer.Close()

	err = writer.WriteHeader([]string{"first_name", "last_name"})
	if err != nil {
		fmt.Println("error writing header: ", err.Error())
		return
	}

	for _, rec := range []ExampleRecord{{"John", "Doe"}, {"Jane", "Doe"}} {
		if err := writer.Write(rec); err != nil {
			fmt.Println("error writing record: ", err.Error())
			break
		}
	}
	// flush our output so that it will be available to print
	writer.Flush()
	fmt.Println(output.String())
	// Output:
	// first_name,last_name
	// John,Doe
	// Jane,Doe
}
