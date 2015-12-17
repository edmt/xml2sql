package main

import (
	"encoding/csv"
	"os"
	"strings"
)

type CSVRecord struct {
	RFC     string
	Id      string
	Ambient string
}

func readCSV(csvFilepath string) map[string]CSVRecord {
	csvfile, _ := os.Open(csvFilepath)
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)

	reader.FieldsPerRecord = -1
	reader.Comma = ','

	rawCSVdata, _ := reader.ReadAll()

	db := make(map[string]CSVRecord)

	var oneRecord CSVRecord
	// var allRecords []CSVRecord

	for _, each := range rawCSVdata {
		oneRecord.RFC = strings.ToUpper(each[0])
		oneRecord.Id = each[1]
		oneRecord.Ambient = strings.ToUpper(each[2])
		//allRecords = append(allRecords, oneRecord)
		db[oneRecord.RFC] = oneRecord
	}

	// return allRecords
	return db
}
