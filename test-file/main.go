package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	//"strings"
)

type LocationsByCountry struct {
	LocID        string
	LocationName string
	Country      string
}

var LocsByCountry []LocationsByCountry

func main() {

	csvFile, err := os.Open("locations-country.txt")

	if err != nil {
		fmt.Println(err)
	}

	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.LazyQuotes = true

	reader.FieldsPerRecord = -1

	csvData, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//strings.Split()
	var oneRecord LocationsByCountry
	var allRecords []LocationsByCountry

	for _, each := range csvData {

		oneRecord.LocID = each[0]
		oneRecord.LocationName = each[1]
		oneRecord.Country = each[2]
		allRecords = append(allRecords, oneRecord)
	}

	jsondata, err := json.Marshal(allRecords) // convert to JSON
	json.Unmarshal(jsondata, &LocsByCountry)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
