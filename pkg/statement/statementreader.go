package statement

import (
	"encoding/csv"
	"log"
	"os"
)

func ReadCsvFile(filepath string) [][]string {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Unable to read input file: "+filepath, err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.FieldsPerRecord = -1
	csvReader.Comma = ';'
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatalf("Unable to parse file as CSV for: "+filepath, err)
	}
	return records
}
