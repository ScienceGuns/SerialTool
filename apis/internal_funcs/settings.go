/*
Copyright © 2025 Mad Scientist Research LLC
This file is part of Serial Tool.
*/

package internal_funcs

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

// CheckForSettingsFile Check if the settings file exists and has the correct mode
func CheckForSettingsFile(settingsFilePath string) bool {
	settingsFile, err := os.Stat(settingsFilePath)
	if err != nil {
		// If the error is not nil, check if this is because the file doesn't exist
		if os.IsNotExist(err) {
			return false
		} else {
			// Print an error but don't return it return false and use CLI
			fmt.Println("Error checking for settings.json. Continuing with CLI args: ", err)
			return false
		}
	}
	// If the mode is correct then we can return true
	if settingsFile.Mode() == 0400 {
		return true
	} else { // Otherwise return false but continue after warning
		fmt.Println("Settings file detected but mode is not 0400!\n" +
			"You SHOULD rotate any credentials in the file after fixing mode.")
		return false
	}
}

// ReadSettingsFile Read the settings.json file and store the data
func ReadSettingsFile(settingsFilePath string) (SerialToolConfig, error) {
	fileData, err := os.ReadFile(settingsFilePath)
	if err != nil {
		return SerialToolConfig{}, fmt.Errorf("failure reading discovered config file: %v", err)
	}

	var outputConfig SerialToolConfig
	err = json.Unmarshal(fileData, &outputConfig)
	if err != nil {
		return SerialToolConfig{}, fmt.Errorf("failure reading discovered config file: %v", err)
	}

	return outputConfig, nil
}

// OutputUserPath output a path under the current user's home directory
func OutputUserPath(userPath string) (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user lookup: ", err)
		return "", fmt.Errorf("failure determining who the current user is: %v", err)
	}
	fullPathOutput := filepath.Join(currentUser.HomeDir, userPath)
	return fullPathOutput, nil
}
