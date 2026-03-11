/*
Copyright © 2025 Mad Scientist Research LLC
This file is part of Serial Tool.
*/

package serials

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// EnsureData checks if the local data files exists, and if not, downloads them.
func EnsureData(dataUrl, localPath string) error { //TODO make this more generic when additional files are needed
	productsFile := filepath.Join(localPath, "products.json")
	// Check if the product file exists
	if _, err := os.Stat(productsFile); os.IsNotExist(err) {
		// Attempt to download it, if it's missing.
		err = DownloadProductData(dataUrl, localPath)
		return err
	}
	return nil
}

// DownloadProductData fetches the remote product data and updates the local file if needed
func DownloadProductData(dataUrl, localPath string) error {
	// Set the local version for comparison
	var localVersion int
	// Set the product data URL and path
	productsUrl := dataUrl + "products.json"
	productsPath := filepath.Join(localPath, "products.json")

	// Try to read the current version. Set to 0 if missing
	fileData, err := os.ReadFile(productsPath)
	if err == nil {
		var localData productDataJSON
		if err := json.Unmarshal(fileData, &localData); err == nil {
			localVersion = localData.Version
		}
	}

	// Fetch the remote product data
	remoteData, err := fetchRemoteData(productsUrl)
	if err != nil {
		return fmt.Errorf("failed to fetch remote products data: %w", err)
	}

	// Compare the two versions using basic math
	if remoteData.Version > localVersion {
		err = saveLocalData(productsPath, remoteData)
		if err != nil {
			return fmt.Errorf("failed to save updated product data: %w", err)
		}
		fmt.Println("The products.json file has been updated!")
		return nil
	}

	// The local version is not lower
	fmt.Println("The products.json file is up to date already.")
	return nil
}

// LoadProductData reads the local products.json file and populates the models into memory
func LoadProductData(localPath string) error {
	var data productDataJSON
	fileData, err := os.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("failed to read local product data: %w", err)
	}

	err = json.Unmarshal(fileData, &data)
	if err != nil {
		return fmt.Errorf("failed to parse local product data: %w", err)
	}

	if data.ProductModel == nil {
		return fmt.Errorf("failed to load any product data")
	}

	productModel = data.ProductModel
	productModelOrder = data.ProductModelOrder
	return nil
}

// fetchRemoteData A helper function to fetch remote data from JSON
func fetchRemoteData(url string) (productDataJSON, error) { //TODO make this more generic when additional files are needed
	var productData productDataJSON
	// Get the product JSON from the web
	response, err := http.Get(url)
	if err != nil {
		return productData, err
	}
	defer response.Body.Close()

	// Fail if not a 200
	if response.StatusCode != http.StatusOK {
		return productData, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	// Return the data and/or any errors
	err = json.NewDecoder(response.Body).Decode(&productData)
	return productData, err
}

// saveLocalData A helper function to save data to a local file as JSON
func saveLocalData(destinationPath string, data productDataJSON) error { //TODO make this more generic when additional files are needed
	// Make sure the data directory exists before trying to save the file
	dataDir := filepath.Dir(destinationPath)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}

	// Marshal the data as JSON
	fileData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	// Write the data to file and return the result as an err
	return os.WriteFile(destinationPath, fileData, 0600)
}
