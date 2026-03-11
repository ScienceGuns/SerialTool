/*
Copyright © 2025 Mad Scientist Research LLC
This file is part of Serial Tool.
*/

package serials

import (
	"fmt"
	"strconv"
)

// Decode decodes an existing USN to a human-readable format
func Decode(serialNumber string) (SerialNumber, error) {
	if len(serialNumber) != 12 {
		return SerialNumber{}, fmt.Errorf("invalid Universal Serial Number length: expected 12, got %d", len(serialNumber))
	}

	modelCode := serialNumber[0:2]
	decodedModel, ok := productModel[modelCode]
	if !ok {
		return SerialNumber{}, fmt.Errorf("unknown model code: %s", modelCode)
	}

	typeCode := rune(serialNumber[2])
	decodedType, ok := productType[typeCode]
	if !ok {
		return SerialNumber{}, fmt.Errorf("unknown type code: %c", typeCode)
	}

	yearCode := serialNumber[3:6]
	yearPrefixVal, err := dateCodeDecoder(rune(yearCode[0]))
	if err != nil {
		return SerialNumber{}, fmt.Errorf("invalid year prefix character: %w", err)
	}
	yearPrefix := (yearPrefixVal - 1) + 20
	yearSuffix, err := strconv.Atoi(yearCode[1:])
	if err != nil {
		return SerialNumber{}, fmt.Errorf("invalid year suffix: %w", err)
	}
	year := yearPrefix*100 + yearSuffix

	supplyChainCode := serialNumber[6:8]

	count, err := strconv.Atoi(serialNumber[8:])
	if err != nil {
		return SerialNumber{}, fmt.Errorf("invalid count segment: %w", err)
	}

	decodedSerial := SerialNumber{
		Model:       decodedModel,
		Type:        decodedType,
		Year:        year,
		SupplyChain: supplyChainCode,
		Count:       count,
	}
	return decodedSerial, nil
}
