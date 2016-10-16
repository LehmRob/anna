package main

import (
	"log"
)

func main() {
	log.Println("Starting WebApp")
	a, err := NewServer(":8080")
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(a.Run())
}
