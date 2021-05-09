package main

import (
	"cdb"
	"log"
)

func main() {
	if err := cdb.GenerateCustomDB(); err != nil {
		log.Fatal(err)
	}
}
