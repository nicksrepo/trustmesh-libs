package main

import (
	 "libcrypto"
		"encoding/hex"
		"fmt"
		"github.com/davecgh/go-spew/spew"
		"log"
)

func main() {

	// Step 1: Define Latitude and Longitude

		lat := 37.8199
		
		lon := -122.4783

		// Step 2: Initialize NetworkAddress
		address, _ := libcrypto.GenerateAddress(lat, lon, 256)

		fmt.Println("Successfully generated a valid NetworkAddress with ZKP.")

		spew.Dump(address)
		log.Println(hex.EncodeToString([]byte(address.PublicKey)))


}
