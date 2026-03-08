package main

import (
	"log"
	"os"
)

func main() {
	log.Println("Crash boot plugin - CRASHING BEFORE ANYTHING!")
	os.Exit(42)
}
