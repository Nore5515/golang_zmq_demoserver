package main

import (
	"fmt"
	"math/rand"
)

// Given an array of strings, return true if string "e" is in the array.
func sliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Create and return two random hex values.
func genRandomPairing() (string, string) {
	x := rand.Int31()
	y := rand.Int31()

	hexX := fmt.Sprintf("%x", x)
	hexY := fmt.Sprintf("%x", y)

	return hexX, hexY
}

// Convert two given integers to hex values, and return.
func getHexPair(x int, y int) (string, string) {
	hexX := fmt.Sprintf("%x", x)
	hexY := fmt.Sprintf("%x", y)

	return hexX, hexY
}
