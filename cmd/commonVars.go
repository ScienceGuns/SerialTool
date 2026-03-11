/*
Copyright © 2025 Mad Scientist Research LLC
This file is part of Serial Tool.
*/

package cmd

import (
	"github.com/ScienceGuns/SerialTool/apis/internal_funcs"
	"runtime"
)

// The version string should be updated before any merge to main
var shortVersion = "0.1.0"
var projectMaintainer = "Mad Scientist Research LLC"
var projectLicense = "MIT"
var functionHelpShort = "Generate/Decode Mad Scientist Research LLC product Universal Serial Numbers"
var functionHelpLong = `This tool is used by Mad Scientist Research LLC to generate or decode serial numbers
using our Universal Serial Number format.`
var buildArch = runtime.GOARCH
var SettingsFilePath, _ = internal_funcs.OutputUserPath(".config/scienceguns/serialtool.json")
var DataFilePath, _ = internal_funcs.OutputUserPath(".config/scienceguns/")
