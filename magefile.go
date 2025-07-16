//go:build mage
// +build mage

package main

import "github.com/magefile/mage/sh"

// Building binaries for the various systems it should run on.
func Build() {
	platforms := map[string]string{
		"linux":   "amd64",
		"darwin":  "arm64",
		"windows": "amd64",
	}
	for os, arch := range platforms {
		suffix := ""
		if os == "windows" {
			suffix = ".exe"
		}
		sh.RunWith(map[string]string{"GOOS": os, "GOARCH": arch},
			"go", "build", "-o", "./build/digitscanlogs-"+os+"-"+arch+suffix, "./cmd/digit_scanlogs")
	}
}

// Cleaning up this repo
func Clean() {
	sh.Run("rm", "-rf", "build")
}

// Copy the binaries to the arcnum server
func Deploy() {
	sh.Run("scp", "build/digitscanlogs-windows*.exe", "rarcnum:")
	sh.Run("scp", "build/digitscanlogs-linux*", "rarcnum:")
}
