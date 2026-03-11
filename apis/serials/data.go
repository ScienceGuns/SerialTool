/*
Copyright © 2025 Mad Scientist Research LLC
This file is part of Serial Tool.
*/

package serials

import "time"

// productModel The map of product models to their serial number encoding
var productModel map[string]string

// productModelOrder defines the order for displaying product models during encoding
var productModelOrder []string

// productDataJSON A struct to keep track of known values in products.json
type productDataJSON struct {
	Version           int               `json:"version"`
	ProductModel      map[string]string `json:"productModel"`
	ProductModelOrder []string          `json:"productModelOrder"`
}

// productType The map of product types to their serial number encoding
var productType = map[rune]string{
	'R': "Rifle",
	'P': "Pistol",
	'F': "Frame",
	'G': "Shotgun",
	'B': "Short Barreled Rifle",
	'H': "Short Barreled Shotgun",
	'S': "Suppressor",
	'A': "Any Other Weapon",
	'M': "Machine Gun",
}

// productTypeOrder defines the order for displaying product types during encoding
var productTypeOrder = []rune{'R', 'P', 'F', 'G', 'B', 'H', 'S', 'A', 'M'}

// serialRecord A USN record in the encoded format for MongoDB
type serialRecord struct {
	SerialNumber string    `bson:"serial_number"`
	CreatedAt    time.Time `bson:"created_at"`
}

// SerialNumber a unique serial number using the Universal Serial Number format
type SerialNumber struct {
	Model       string
	Type        string
	Year        int
	SupplyChain string
	Count       int
}
