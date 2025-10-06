//go:build mage
// +build mage

package main

import "github.com/magefile/mage/sh"

var platforms = map[string]string{
	"linux":   "amd64",
	"darwin":  "arm64",
	"windows": "amd64",
}

// Building binaries for the various systems it should run on.
func Build() {
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
	for os, arch := range platforms {
		if os == "darwin" {
			continue
		}
		suffix := ""
		if os == "windows" {
			suffix = ".exe"
		}
		sh.Run("scp", "build/digitscanlogs-"+os+"-"+arch+suffix, "rarcnum:")
	}
}
