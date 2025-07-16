package main

import (
	"context"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/edipretoro/digit_scanlogs/internal/database"
	"github.com/google/uuid"
)

func digestFile(dir string, file string) (string, error) {
	f, err := os.Open(filepath.Join(dir, file))
	if err != nil {
		return "", fmt.Errorf("opening file for digesting: %s", err)
	}
	defer f.Close()

	h := sha512.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("copying data for digesting: %s", err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func IsDigitProject(projectPath string) bool {
	path, err := filepath.Abs(projectPath)
	if err != nil {
		log.Println("Error getting absolute path:", err)
		return false
	}
	tiffFiles, err := filepath.Glob(filepath.Join(path, "*.tif"))
	if err != nil {
		log.Println("Error checking for TIFF files:", err)
		return false
	}
	return len(tiffFiles) > 0
}

func processProject(project database.Project, dbQueries *database.Queries) error {
	fmt.Printf("Processing project at: %s\n", project.Name)
	return nil
}

func processingScanDirectory(scanDir string, dbQueries *database.Queries) error {
	projects, err := os.ReadDir(scanDir)
	if err != nil {
		return fmt.Errorf("failed to read scan directory: %w", err)
	}
	for _, project := range projects {
		if project.IsDir() {
			projectPath := filepath.Join(scanDir, project.Name())
			if IsDigitProject(projectPath) {
				userDb, err := GetUserFromPath(projectPath, dbQueries)
				if err != nil {
					return fmt.Errorf("failed to get user from project path %s: %w", projectPath, err)
				}
				var projectDb database.Project
				projectDb, err = dbQueries.GetProjectByPath(context.Background(), projectPath)
				if err != nil {
					if err == sql.ErrNoRows {
						projectDb, err = dbQueries.CreateProject(
							context.Background(),
							database.CreateProjectParams{
								ID:          uuid.NewString(),
								Name:        project.Name(),
								Path:        projectPath,
								Description: sql.NullString{},
								CreatedBy:   userDb.Uid,
							},
						)
						if err != nil {
							return fmt.Errorf("failed to create project in database: %w", err)
						}
					} else {
						return fmt.Errorf("failed to get project by path %s: %w", projectPath, err)
					}
				}
				err = processProject(projectDb, dbQueries)
				if err != nil {
					return fmt.Errorf("failed to process project %s: %w", project.Name(), err)
				}
			}
		}
	}
	return nil
}
