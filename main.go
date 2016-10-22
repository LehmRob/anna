package main

import (
	"log"
)

func main() {
	log.Println("Starting WebApp")
	a, err := NewServer(":443")
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(a.Run())
}
