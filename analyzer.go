package main

import "fmt"

// Analyzer defines functions for diffrent analyzer implementations
type Analyzer interface {
	Analyze() ([]Result, error)
}

// Result holds the data for a line in the data table
type Result struct {
	// Name is the name of the person
	Name string

	// Address contains the street and the house number
	Address string

	// ZipCode is the zip code for the city the person live
	ZipCode string

	// City is the name of the city where the person live
	City string
}

func (r *Result) String() string {
	return fmt.Sprintf("%s; %s %s %s", r.Name, r.Address, r.ZipCode, r.City)
}
