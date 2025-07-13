package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"github.com/edipretoro/digit_scanlogs/internal/database"
	"github.com/google/uuid"
)

func IsDigitProject(projectPath string) bool {
	path, err := filepath.Abs(projectPath)
	if err != nil {
		log.Println("Error getting absolute path:", err)
		return false
	}
	tiffFiles, err := filepath.Glob(fmt.Sprintf("%s/*.tif", path))
	if err != nil {
		log.Println("Error checking for TIFF files:", err)
		return false
	}
	return len(tiffFiles) > 0
}

func processProject(projectPath string) error {
	fmt.Printf("Processing project at: %s\n", projectPath)
	return nil
}

func processingScanDirectory(scanDir string, dbQueries *database.Queries) error {
	projects, err := os.ReadDir(scanDir)
	if err != nil {
		return fmt.Errorf("failed to read scan directory: %w", err)
	}
	for _, project := range projects {
		if project.IsDir() {
			projectPath := fmt.Sprintf("%s/%s", scanDir, project.Name())
			if IsDigitProject(projectPath) {
				fmt.Printf("Found Digit project: %s\n", project.Name())
				err = processProject(projectPath)
				if err != nil {
					return fmt.Errorf("failed to process project %s: %w", project.Name(), err)
				}
			}
		}
	}
	return nil
}
