package main

import (
	"log"
	"os"
)

//go:generate ./package.nu

func main() {
	log.Println("Crash boot plugin - CRASHING BEFORE ANYTHING!")
	os.Exit(42)
}
