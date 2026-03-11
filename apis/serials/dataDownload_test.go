/*
Copyright © 2025 Mad Scientist Research LLC
This file is part of Serial Tool.
*/

package serials

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// TestSaveAndLoadProductData tests reading and writing the JSON to the local filesystem
func TestSaveAndLoadProductData(t *testing.T) {
	// Create a temporary directory for any test files
	tempDir := t.TempDir()
	testFilePath := filepath.Join(tempDir, "products.json")

	// Create mock product data to save
	mockData := productDataJSON{
		Version:           1,
		ProductModel:      map[string]string{"ZZ": "Test Product DO NOT USE"},
		ProductModelOrder: []string{"ZZ"},
	}

	// Test saving the file locally
	err := saveLocalData(testFilePath, mockData)
	if err != nil {
		t.Fatalf("saveLocalData failed: %v", err)
	}

	// Verify the file was actually saved
	if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
		t.Fatalf("Expected file to be created at %s, but it does not exist", testFilePath)
	}

	// Test loading the local file
	err = LoadProductData(testFilePath)
	if err != nil {
		t.Fatalf("LoadProductData failed: %v", err)
	}

	// Verify the variables get populated correctly via this file
	if productModel["ZZ"] != "Test Product DO NOT USE" {
		t.Errorf("Expected productModel to contain 'ZZ':'Test Product DO NOT USE', got %v", productModel)
	}
	if !reflect.DeepEqual(productModelOrder, []string{"ZZ"}) {
		t.Errorf("Expected productModelOrder to be ['ZZ'], got %v", productModelOrder)
	}
}

// TestDownloadProductData tests the version comparison and downloading logic
func TestDownloadProductData(t *testing.T) {
	// Create a temporary directory for any test files
	tempDir := t.TempDir()

	// Control what version our mock server responds with
	mockRemoteVersion := 202506051000

	// Create a mock HTTP server to avoid actually pulling from the internet
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := productDataJSON{
			Version:           mockRemoteVersion,
			ProductModel:      map[string]string{"ZZ": "Test Product DO NOT USE"},
			ProductModelOrder: []string{"ZZ"},
		}
		// Write the mock data
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			return
		}
	}))
	defer testServer.Close()

	// The download function appends "products.json" to the dataUrl, so ensure it ends with a slash
	dataUrl := testServer.URL + "/"

	// Initial a download as if there is no local file
	err := DownloadProductData(dataUrl, tempDir)
	if err != nil {
		t.Fatalf("Initial DownloadProductData failed: %v", err)
	}

	// Load the data to verify it pulled version 202506051000
	productsPath := filepath.Join(tempDir, "products.json")
	err = LoadProductData(productsPath)
	if err != nil {
		t.Fatalf("Failed to load freshly downloaded data: %v", err)
	}
	if productModel["ZZ"] != "Test Product DO NOT USE" {
		t.Errorf("Expected Test Product DO NOT USE to be loaded into memory")
	}

	// Test version comparison with an older remote
	mockRemoteVersion = 202501011000 // Lower the remote version to 202501011000
	err = DownloadProductData(dataUrl, tempDir)
	if err != nil {
		t.Fatalf("DownloadProductData failed when remote was older: %v", err)
	}

	// Read the file directly to make sure it wasn't overwritten by the older version
	fileData, _ := os.ReadFile(productsPath)
	var localData productDataJSON
	err = json.Unmarshal(fileData, &localData)
	if err != nil {
		return
	}
	if localData.Version != 202506051000 {
		t.Errorf("Expected local version to remain 202506051000, but it changed to %d", localData.Version)
	}

	// Test version comparison with a newer remote
	mockRemoteVersion = 202603111200
	err = DownloadProductData(dataUrl, tempDir)
	if err != nil {
		t.Fatalf("DownloadProductData failed when remote was newer: %v", err)
	}

	// Verify the file was updated to version 202603111200
	fileData, _ = os.ReadFile(productsPath)
	err = json.Unmarshal(fileData, &localData)
	if err != nil {
		return
	}
	if localData.Version != 202603111200 {
		t.Errorf("Expected local version to update to 202603111200, but it is %d", localData.Version)
	}
}

// TestEnsureData tests that files are only downloaded if they don't exist
func TestEnsureData(t *testing.T) {
	// Create a temporary directory for any test files
	tempDir := t.TempDir()

	// Create a mock HTTP server to avoid actually pulling from the internet
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"version": 202506051000, "productModel": {"ZZ":"Test Product DO NOT USE"}, "productModelOrder": ["ZZ"]}`))
		if err != nil {
			return
		}
	}))
	defer testServer.Close()

	// The download function appends "products.json" to the dataUrl, so ensure it ends with a slash
	dataUrl := testServer.URL + "/"

	// File doesn't exist and should be downloaded
	err := EnsureData(dataUrl, tempDir)
	if err != nil {
		t.Fatalf("EnsureData failed on missing file: %v", err)
	}

	productsPath := filepath.Join(tempDir, "products.json")
	if _, err := os.Stat(productsPath); os.IsNotExist(err) {
		t.Fatalf("EnsureData did not create the file")
	}

	//The file exists andEnsureData should do nothing
	// An invalid URL is given to make sure it fails if it tries to download anyway
	invalidDataUrl := "http://127.0.0.1:0/invalid/"
	err = EnsureData(invalidDataUrl, tempDir)
	if err != nil {
		t.Fatalf("EnsureData attempted to download despite file existing, resulting in error: %v", err)
	}
}
