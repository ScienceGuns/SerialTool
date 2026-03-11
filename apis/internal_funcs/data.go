/*
Copyright © 2025 Mad Scientist Research LLC
This file is part of Serial Tool.
*/

package internal_funcs

// SerialToolConfig A struct to keep track of known values in settings.json
type SerialToolConfig struct {
	DataURL string `json:"data_url"`
	MongoDB struct {
		URI        string `json:"uri"`
		Database   string `json:"database"`
		Collection string `json:"collection"`
	} `json:"mongodb"`
}
