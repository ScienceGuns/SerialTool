/*
Copyright © 2025 Mad Scientist Research LLC
This file is part of Serial Tool.
*/

package serials

import (
	"testing"
)

// TestDecode test the Decode function for valid and invalid cases
func TestDecode(t *testing.T) {
	// Map some test product.json data since we aren't loading files
	productModel = map[string]string{
		"TI": "The Isotope",
		"TG": "The Gauss",
		"TK": "The Kinetic",
		"TR": "The Repulsor",
		"TV": "The Vector",
		"TT": "The Tensor",
		"TC": "The Catalyst",
		"ZZ": "Test Product DO NOT USE",
	}
	productModelOrder = []string{"TI", "TK", "TG", "TR", "TV", "TT", "TC", "ZZ"}

	// Define test cases
	testCases := []struct {
		name         string
		serialNumber string
		expectError  bool
	}{
		{
			name:         "Valid Universal Serial Number",
			serialNumber: "TIR125MS0123", // A valid "The Isotope" Rifle made in 2025 as count number 123
			expectError:  false,
		},
		{
			name:         "Valid Variance Marking Serial",
			serialNumber: "TIB125SG0001", // A valid "The Isotope" SBR made in 2025
			expectError:  false,
		},
		{
			name:         "Invalid Length - USN Too Short",
			serialNumber: "TKR125MS123", // An otherwise valid serial but missing the padded 0 for the count
			expectError:  true,
		},
		{
			name:         "Invalid Length - USN Too Long",
			serialNumber: "TCR125MSS0123", // Example of someone using double digits for the day vs a one character code
			expectError:  true,
		},
		{
			name:         "Invalid Year Prefix",
			serialNumber: "TIR-24MS0123",
			expectError:  true,
		},
		{
			name:         "Invalid Year Suffix",
			serialNumber: "TIRAXXMS0123",
			expectError:  true,
		},
		{
			name:         "Invalid Count Segment",
			serialNumber: "TIR125MSABCD",
			expectError:  true,
		},
	}

	// Loop through test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Decode(tc.serialNumber)

			// Check if we got an error when we didn't expect one
			if !tc.expectError && err != nil {
				t.Errorf("Decode() with serial '%s' returned an unexpected error: %v", tc.serialNumber, err)
			}

			// Check if we didn't get an error when we expected one
			if tc.expectError && err == nil {
				t.Errorf("Decode() with serial '%s' was expected to return an error, but it did not", tc.serialNumber)
			}
		})
	}
}
