package main

import (
	"encoding/csv"
	"io"
	"os"
	"regexp"
	"sort"
)

const (
	iName      = 0
	iAddress   = 1
	iPlace     = 2
	iBirthdate = 3
	iAdmission = 4
	iDischarge = 5
	iCategory  = 6
	iDuration  = 7
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
func (a *CsvAnalyzer) Analyze() (ZipCodeGroups, error) {
	childreen, err := a.parseFile()
	if err != nil {
		return nil, err
	}

	return toZipCodeGroups(childreen), nil
}

func toZipCodeGroups(childreen Childreen) []ZipCodeGroup {
	zipCodes := parseZipCodes(childreen)

	var groups []ZipCodeGroup
	for _, zipCode := range zipCodes {
		groups = append(groups, createZipCodeGroup(childreen, zipCode))
	}

	return groups
}

func parseZipCodes(childreen []Child) []string {
	var zipCodes []string
	for _, child := range childreen {
		if !isZipCodeAvail(zipCodes, child.ZipCode) {
			zipCodes = append(zipCodes, child.ZipCode)
		}
	}

	return zipCodes
}

func isZipCodeAvail(zipCodes []string, askedZipCode string) bool {
	if len(zipCodes) == 0 {
		return false
	}

	for _, zipCode := range zipCodes {
		if zipCode == askedZipCode {
			return true
		}
	}

	return false
}

func createZipCodeGroup(childreen []Child, zipCode string) ZipCodeGroup {
	var zipCodeChildreen []Child
	for _, child := range childreen {
		if child.ZipCode == zipCode {
			zipCodeChildreen = append(zipCodeChildreen, child)
		}
	}

	zipCodeGroup := ZipCodeGroup{
		ZipCode:   zipCode,
		Childreen: zipCodeChildreen,
	}

	zipCodeGroup.sort()
	return zipCodeGroup
}

func (a *CsvAnalyzer) parseFile() (Childreen, error) {
	var childreen Childreen

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

		child, err := createChild(record)
		if err != nil {
			return nil, err
		}
		childreen = append(childreen, child)
	}
	sort.Sort(childreen)
	return childreen, nil
}

func createChild(record []string) (Child, error) {
	zipCodeRegex, err := regexp.Compile("\\d{5}")
	if err != nil {
		return Child{}, err
	}

	cityRegex, err := regexp.Compile("[A-z]+")
	if err != nil {
		return Child{}, err
	}

	zipCode := zipCodeRegex.FindString(record[iPlace])
	city := cityRegex.FindString(record[iPlace])

	return Child{
		Name:          record[iName],
		Address:       record[iAddress],
		ZipCode:       zipCode,
		City:          city,
		Birthdate:     record[iBirthdate],
		Category:      record[iCategory],
		AdmissionDate: record[iAdmission],
		DischargeDate: record[iDischarge],
		Duration:      record[iDuration],
	}, nil
}
