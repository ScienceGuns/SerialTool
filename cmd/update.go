/*
Copyright © 2025 Mad Scientist Research LLC
This file is part of Serial Tool.
*/

package cmd

import (
	"fmt"
	"github.com/ScienceGuns/SerialTool/apis/internal_funcs"
	"github.com/ScienceGuns/SerialTool/apis/serials"
	"github.com/spf13/cobra"
	"os"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update SerialTool's files.",
	Long:  `Pull the latest products and other files from Mad Scientist Research LLC.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !internal_funcs.CheckForSettingsFile(SettingsFilePath) {
			os.Exit(1)
		}
		// Pull the settings in
		settings, err := internal_funcs.ReadSettingsFile(SettingsFilePath)
		if err != nil {
			fmt.Println(err)
		}
		if settings.DataURL == "" {
			fmt.Println("Error: 'data_url' is not configured in settings.json.")
			os.Exit(1)
		}

		err = serials.DownloadProductData(settings.DataURL, DataFilePath)
		if err != nil {
			fmt.Printf("Error updating product data: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
