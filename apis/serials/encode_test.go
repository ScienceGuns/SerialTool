/*
Copyright © 2025 Mad Scientist Research LLC
This file is part of Serial Tool.
*/

package serials

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/ScienceGuns/SerialTool/apis/internal_funcs"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

// mockDatastore mock connecting to a MongoDB database instance
type mockDatastore struct {
	countToReturn int
	shouldError   bool
	config        internal_funcs.SerialToolConfig
}

// GetSerialCount mock getting a specified count of existing USNs
func (m *mockDatastore) GetSerialCount(config internal_funcs.SerialToolConfig, prefix string) (int, error) {
	if m.shouldError {
		return 0, fmt.Errorf("mock database error")
	}
	return m.countToReturn, nil
}

// StoreSerialNumber mock storing a USN with a happy response
func (m *mockDatastore) StoreSerialNumber(config internal_funcs.SerialToolConfig, serial string) error {
	if m.shouldError {
		return fmt.Errorf("mock database error")
	}
	return nil
}

// GetConfig mock getting a config from settings
func (m *mockDatastore) GetConfig() internal_funcs.SerialToolConfig {
	return m.config
}

// Disconnect mock the database disconnect.
func (m *mockDatastore) Disconnect() {}

// TestEncode test the Encode function with a single basic test for now.
func TestEncode(t *testing.T) {
	// Set a fixed point in time for testing
	fixedPointInTime := time.Date(2025, time.June, 5, 10, 0, 0, 0, time.UTC)
	// Save the real time for later
	realTime := now
	// Override the now variable with our fixed point in time
	now = func() time.Time {
		return fixedPointInTime
	}
	// Return to original time after the test so we don't break any other tests
	defer func() {
		now = realTime
	}()

	input := "TI\nF\nMS\n\n"                            // Specify a single "The Isotope" as a Frame configuration with the MS supply chain code
	reader := bufio.NewReader(strings.NewReader(input)) // Push those inputs into the buffer automatically

	// Mock our DB connection and return count
	mockDB := &mockDatastore{
		countToReturn: 22, // Lie about 22 existing Isotopes being created today
		shouldError:   false,
		config:        internal_funcs.SerialToolConfig{},
	}

	// Capture our standard input
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	//Pass the mock database and data to the USN encoder
	Encode(mockDB, reader)

	// Cleanup of resources
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("Failed to read test output: %v", err)
	}
	output := buf.String()

	// Check our assertions for validity
	if !strings.Contains(output, "Generated Universal Serial Number: TIF125MS0023") {
		t.Errorf("Expected USN not found in output: '%s'", output)
	}
	if !strings.Contains(output, "Successfully saved to database.") {
		t.Errorf("Expected success message not found in output: '%s'", output)
	}
}

// TestEncode_Variance test the Encode function with a new variance marked receiver
func TestEncode_Variance(t *testing.T) {
	// Set a fixed point in time for testing
	fixedPointInTime := time.Date(2025, time.June, 5, 10, 0, 0, 0, time.UTC)
	// Save the real time for later
	realTime := now
	// Override the now variable with our fixed point in time
	now = func() time.Time {
		return fixedPointInTime
	}
	// Return to original time after the test so we don't break any other tests
	defer func() {
		now = realTime
	}()

	input := "TV\nR\nSG\n1\n"                           // Specify a single "The Vector" as a Rifle configuration with the SG supply chain code
	reader := bufio.NewReader(strings.NewReader(input)) // Push those inputs into the buffer automatically

	// Mock our DB connection and return count
	mockDB := &mockDatastore{
		countToReturn: 5, // Lie about 5 existing Vectors being created today
		shouldError:   false,
		config:        internal_funcs.SerialToolConfig{},
	}

	// Capture our standard input
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	//Pass the mock database and data to the USN encoder
	Encode(mockDB, reader)

	// Cleanup of resources
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("Failed to read test output: %v", err)
	}
	output := buf.String()

	// Check our assertions for validity
	if !strings.Contains(output, "Generated Universal Serial Number: TVR125SG0006") {
		t.Errorf("Expected variance marking USN not found in output: '%s'", output)
	}
	if !strings.Contains(output, "Successfully saved to database.") {
		t.Errorf("Expected success message not found in output: '%s'", output)
	}
}

// TestEncode_Batch test the Encode function by making 3 variance marked receivers
func TestEncode_Batch(t *testing.T) {
	// Set a fixed point in time for testing
	fixedPointInTime := time.Date(2025, time.June, 5, 10, 0, 0, 0, time.UTC)
	// Save the real time for later
	realTime := now
	// Override the now variable with our fixed point in time
	now = func() time.Time {
		return fixedPointInTime
	}
	// Return to original time after the test so we don't break any other tests
	defer func() {
		now = realTime
	}()

	input := "TI\nB\nMS\n3\n"                           // Specify 3 "The Isotope" as an SBR configuration with the MS supply chain code
	reader := bufio.NewReader(strings.NewReader(input)) // Push those inputs into the buffer automatically

	// Mock our DB connection and return count
	mockDB := &mockDatastore{
		countToReturn: 0, // Lie that 0 previous models have been created
		shouldError:   false,
		config:        internal_funcs.SerialToolConfig{},
	}

	// Capture our standard input
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	//Pass the mock database and data to the USN encoder
	Encode(mockDB, reader)

	// Cleanup of resources
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("Failed to read test output: %v", err)
	}
	output := buf.String()

	// Check our assertions for validity
	if !strings.Contains(output, "Generated Universal Serial Number: TIB125MS0001") {
		t.Errorf("Expected first variance marking USN not found in output: '%s'", output)
	}
	if !strings.Contains(output, "Generated Universal Serial Number: TIB125MS0002") {
		t.Errorf("Expected second variance marking USN not found in output: '%s'", output)
	}
	if !strings.Contains(output, "Generated Universal Serial Number: TIB125MS0003") {
		t.Errorf("Expected third variance marking USN not found in output: '%s'", output)
	}
	if strings.Count(output, "Successfully saved to database.") != 3 {
		t.Errorf("Expected 3 database success messages to be found in output: '%s'", output)
	}
}
