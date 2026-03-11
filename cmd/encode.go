/*
Copyright © 2025 Mad Scientist Research LLC
This file is part of Serial Tool.
*/

package cmd

import (
	"bufio"
	"fmt"
	"github.com/ScienceGuns/SerialTool/apis/internal_funcs"
	"github.com/ScienceGuns/SerialTool/apis/serials"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// hostnameCmd represents the hostname command
var encodeCmd = &cobra.Command{
	Use:   "encode",
	Short: "Encode a new USN",
	Long:  `Encode a new Universal Serial Number for a product manufactured today.`,
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
		// Connect to the database
		db, err := serials.InitDB(settings)
		if err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}
		// Prep for Database disconnect when done
		defer db.Disconnect()
		// Follow the guided serialization process
		reader := bufio.NewReader(os.Stdin)
		serials.Encode(db, reader)
	},
}

func init() {
	rootCmd.AddCommand(encodeCmd)
}
