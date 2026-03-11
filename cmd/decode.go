/*
Copyright © 2025 Mad Scientist Research LLC
This file is part of Serial Tool.
*/

package cmd

import (
	"fmt"
	"github.com/ScienceGuns/SerialTool/apis/internal_funcs"
	"github.com/ScienceGuns/SerialTool/apis/serials"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// decodeCmd represents the hostname command
var decodeCmd = &cobra.Command{
	Use:   "decode <USN>",
	Short: "Decode an existing USN",
	Long:  `Decode a Universal Serial Number from a previously manufactured product.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Pull the settings in
		if !internal_funcs.CheckForSettingsFile(SettingsFilePath) {
			os.Exit(1)
		}
		settings, err := internal_funcs.ReadSettingsFile(SettingsFilePath)
		if err != nil {
			fmt.Println(err)
		}
		// Validate product data has been downloaded
		if err := serials.EnsureData(settings.DataURL, DataFilePath); err != nil {
			fmt.Printf("Error validating product data exists: %v\n", err)
			os.Exit(1)
		}
		// Load the product data into memory
		productsFile := filepath.Join(DataFilePath, "products.json")
		if err := serials.LoadProductData(productsFile); err != nil {
			fmt.Printf("Error loading product data: %v\n", err)
			os.Exit(1)
		}
		serialNumber := args[0]
		decodedSerial, err := serials.Decode(serialNumber)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		returnMessage := fmt.Sprintf("This USN decodes to:\n"+
			"Model: %s\n"+
			"Type At Manufacture: %s\n"+
			"Year Made: %d\n"+
			"Supply Chain Code: %s\n"+
			"Year's Count Number: %d\n"+
			"Note: Supply chain codes are tracked internally.",
			decodedSerial.Model, decodedSerial.Type, decodedSerial.Year, decodedSerial.SupplyChain, decodedSerial.Count)
		fmt.Println(returnMessage)
	},
}

func init() {
	rootCmd.AddCommand(decodeCmd)
}
