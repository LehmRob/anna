package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	catKr = "KR"
	catKg = "KG"
	catH  = "H"
)

var header = []string{
	"Name", "Adresse", "Ort", "Geburtsdatum", "Eintritt", "Austritt", "Einrichtung", "Betreuungszeit",
}

// ZipCodeGroup holds together all childs with the same ZipCode
type ZipCodeGroup struct {
	ZipCode   string
	Childreen Childreen
}

// Child structure which contains all availible date to a child
type Child struct {
	// Name is the name of the person
	Name string
	// Address contains the street and the house number
	Address string
	// ZipCode is the zip code for the city the person live
	ZipCode string
	// City is the name of the city where the person live
	City string
	// Birthdate
	Birthdate string
	// Institution
	Category string
	// Date of the admission of the child
	AdmissionDate string
	// Date of the discharge; this could empty
	DischargeDate string
	// Duration is the care duration per day of the child in hours
	Duration string
}

// ZipCodeGroups is a array of ZipCodeGroups
type ZipCodeGroups []ZipCodeGroup

// Childreen is a array of Child elements
type Childreen []Child

func (c *Child) toArray() []string {
	return []string{c.Name, c.Address, fmt.Sprintf("%s, %s", c.ZipCode, c.City),
		c.Birthdate, c.AdmissionDate, c.DischargeDate, c.Category, c.Duration}
}

// Len returns the count of childreen
func (c Childreen) Len() int {
	return len(c)
}

// Less returns true if the element on the index i is lesser than the element at index
// j. Else return false
func (c Childreen) Less(i, j int) bool {
	return c[i].Name < c[j].Name
}

// Swap swaps the element on index i with the element on index j
func (c Childreen) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (z *ZipCodeGroup) writeCsv() error {
	f, err := os.OpenFile(filepath.Join("/tmp", z.ZipCode+".csv"),
		os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	csvFile := csv.NewWriter(f)
	err = csvFile.Write(header)
	if err != nil {
		return err
	}

	for _, child := range z.Childreen {
		err := csvFile.Write(child.toArray())
		if err != nil {
			return err
		}

		log.Println(child.toArray())
	}

	csvFile.Flush()

	return nil
}

func (z *ZipCodeGroup) sort() {
	kr := z.extractCat(catKr)
	kg := z.extractCat(catKg)
	h := z.extractCat(catH)

	z.Childreen = merge(kr, kg, h)
}

func (z *ZipCodeGroup) extractCat(cat string) []Child {
	var group []Child
	for _, child := range z.Childreen {
		if child.Category == cat {
			group = append(group, child)
		}
	}
	return group
}

func merge(krGroup, kgGroup, hGroup []Child) []Child {
	var group []Child

	for _, kr := range krGroup {
		group = append(group, kr)
	}

	for _, kg := range kgGroup {
		group = append(group, kg)
	}

	for _, h := range hGroup {
		group = append(group, h)
	}

	return group
}
