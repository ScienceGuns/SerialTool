/*
Copyright © 2025 Mad Scientist Research LLC
This file is part of Serial Tool.
*/

package serials

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// now is a variable that can be used or overridden for testing via patching
var now = time.Now

// Encode guide the user through encoding a newly manufactured product through serialization with the Universal Serial Number format
func Encode(db Datastore, reader *bufio.Reader) {
	// Get the config from the interface
	config := db.GetConfig()
	fmt.Printf("\n--- Encode New Serial Number ---\n\n")

	// Display a list of models to choose from using an ordered slice
	fmt.Println("Select a Model:")
	for _, code := range productModelOrder {
		// Look up the name from the map using the forced order
		name := productModel[code]
		fmt.Printf("  %s: %s\n", code, name)
	}
	// Only permit a valid selection before continuing
	var modelCode string
	for {
		fmt.Print("Enter model code: ")
		input, _ := reader.ReadString('\n')
		modelCode = strings.TrimSpace(strings.ToUpper(input))
		if _, exists := productModel[modelCode]; exists {
			break
		}
		fmt.Println("Invalid model code. Please try again.")
	}

	// Display a list of types to choose from
	fmt.Println("\nSelect a Product Type:")
	for _, code := range productTypeOrder {
		name := productType[code]
		fmt.Printf("  %c: %s\n", code, name)
	}
	// Only permit a valid selection before continuing
	var typeCode rune
	for {
		fmt.Print("Enter type code: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToUpper(input))
		if len(input) == 1 {
			typeCode = rune(input[0])
			if _, exists := productType[typeCode]; exists {
				break
			}
		}
		fmt.Println("Invalid type code. Please try again.")
	}

	// Ask about supply chain code
	var supplyChainCode string
	for {
		fmt.Print("\nEnter 2-character supply chain code (Default: MS): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToUpper(input))
		if input == "" {
			supplyChainCode = "MS"
			break
		}
		if len(input) == 2 {
			supplyChainCode = input
			break
		}
		fmt.Println("Invalid supply chain code. Please enter 2 characters.")
	}

	// Ask how many units will be created
	var generationCount int
	for {
		fmt.Print("\nHow many serials of this product should be generated? (Default: 1): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			generationCount = 1
			break
		}
		num, err := strconv.Atoi(input)
		if err != nil || num < 1 {
			fmt.Println("Invalid count. Please enter a positive number or press [Enter] for 1.")
			continue
		}
		generationCount = num
		break
	}

	// Get the current date values and convert it to the Universal Serial Number format
	currentTime := now()
	year, _, _ := currentTime.Date()

	// Encode the year into the USN format
	encodedYear, err := yearEncoder(year)
	if err != nil {
		log.Printf("Error encoding year: %v\n", err)
		return
	}

	// Store the model/type/year/supplyChainCode prefix as a variable
	prefix := fmt.Sprintf("%s%c%s%s", modelCode, typeCode, encodedYear, supplyChainCode)

	// Get a count of current units for this prefix
	count, err := db.GetSerialCount(config, prefix)
	if err != nil {
		log.Printf("Error getting count from database: %v\n", err)
		return
	}

	// Loop through the number of weapons to produce
	for weaponCount := 0; weaponCount < generationCount; weaponCount++ {
		// Increment that count for the new weapon without re-reading the database each time
		nextCount := count + 1 + weaponCount
		countStr := fmt.Sprintf("%04d", nextCount)
		// Store and display the final USN for the new weapon
		finalSerial := prefix + countStr
		fmt.Printf("\nGenerated Universal Serial Number: %s\n", finalSerial)

		// Save the new USN to the database
		if err := db.StoreSerialNumber(config, finalSerial); err != nil {
			log.Printf("Error saving USN to database: %v\n", err)
			log.Printf("DO NOT PRODUCE THE WEAPON UNTIL THIS IS CORRECTED!")
			return
		} else {
			fmt.Println("Successfully saved to database.")
		}
	}
}
