/*
Copyright © 2025 Mad Scientist Research LLC
This file is part of Serial Tool.
*/

package serials

import "fmt"

// dateCodeDecoder convert letters to numbers for use with the date portion of the serial string
func dateCodeDecoder(dateCode rune) (int, error) {
	// Do not decode if it's a single character number
	if dateCode >= '1' && dateCode <= '9' {
		return int(dateCode - '0'), nil
	}
	// Decode letters to a number starting with A as 10
	if dateCode >= 'A' && dateCode <= 'Z' {
		return int(dateCode-'A') + 10, nil
	}
	return 0, fmt.Errorf("invalid character for conversion: %c", dateCode)
}

// Convert numbers to letters (rune not string) for use with the date portion of the serial string
func dateCodeEncoder(dateValue int) (rune, error) {
	// Do not bother encoding single digit numbers as letters
	if dateValue >= 1 && dateValue <= 9 {
		return rune('0' + dateValue), nil
	}
	// Encode 10+ as a single character letter
	if dateValue >= 10 && dateValue <= 35 {
		return rune('A' + dateValue - 10), nil
	}
	return ' ', fmt.Errorf("value out of range for conversion: %d", dateValue)
}

// yearEncoder converts a 4-digit year into the 3-character format.
func yearEncoder(year int) (string, error) {
	if year < 2000 || year > 3599 {
		return "", fmt.Errorf("year %d is out of range (2000-3599)", year)
	}
	prefixVal := (year/100 - 20) + 1
	prefixChar, err := dateCodeEncoder(prefixVal)
	if err != nil {
		return "", err
	}
	suffix := fmt.Sprintf("%02d", year%100)
	return string(prefixChar) + suffix, nil
}
