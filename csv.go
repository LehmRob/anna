package main

import (
	"encoding/csv"
	"io"
	"os"
	"regexp"
)

const (
	iName    = 0
	iAddress = 1
	iPlace   = 2
)

// CsvAnalyzer is the object for analyzing the csv data.
type CsvAnalyzer struct {
	path string
}

// NewCsvAnalyzer creates a new instance of CsvAnalyzer.
// The function needs the path for the file and returns a pointer
// to the new created.
func NewCsvAnalyzer(filePath string) *CsvAnalyzer {
	a := new(CsvAnalyzer)
	a.path = filePath

	return a
}

// Analyze parses the data from the csv file and returns the sorted results
// as array of Result structure
func (a *CsvAnalyzer) Analyze() ([]Result, error) {
	var results []Result

	f, err := os.Open(a.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)

	firstLine := true
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		if firstLine {
			firstLine = false
			continue
		}

		resultElem, err := createResult(record)
		if err != nil {
			return nil, err
		}
		results = append(results, resultElem)
	}

	return results, nil
}

func createResult(record []string) (Result, error) {
	zipCodeRegex, err := regexp.Compile("\\d{5}")
	if err != nil {
		return Result{}, err
	}

	cityRegex, err := regexp.Compile("[A-z]+")
	if err != nil {
		return Result{}, err
	}

	zipCode := zipCodeRegex.FindString(record[iPlace])
	city := cityRegex.FindString(record[iPlace])

	return Result{
		Name:    record[iName],
		Address: record[iAddress],
		ZipCode: zipCode,
		City:    city,
	}, nil
}
