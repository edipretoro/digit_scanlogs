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
	"path/filepath"
	"sync"

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
	files, err := os.ReadDir(project.Path)
	if err != nil {
		return fmt.Errorf("failed to read project directory: %w", err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue // Skip directories for now
		}
		metadata, err := os.Stat(filepath.Join(project.Path, file.Name()))
		if err != nil {
			return fmt.Errorf("failed to get metadata for file %s: %w", file.Name(), err)
		}
		var fileDb database.File
		sha512Digest, err := digestFile(project.Path, file.Name())
		if err != nil {
			return fmt.Errorf("failed to digest file %s: %w", file.Name(), err)
		}
		userDb, err := GetUserFromPath(filepath.Join(project.Path, file.Name()), dbQueries)
		if err != nil {
			return fmt.Errorf("failed to get user from file path %s: %w", filepath.Join(project.Path, file.Name()), err)
		}
		fileDb, err = dbQueries.GetFileByPath(context.Background(), filepath.Join(project.Path, file.Name()))
		if err != nil {
			if err == sql.ErrNoRows {
				fileDb, err = dbQueries.CreateFile(
					context.Background(),
					database.CreateFileParams{
						ID:          uuid.NewString(),
						Name:        file.Name(),
						ProjectID:   project.ID,
						UserID:      userDb.Uid,
						Path:        filepath.Join(project.Path, file.Name()),
						Size:        metadata.Size(),
						Mode:        metadata.Mode().String(),
						Modtime:     metadata.ModTime(),
						Sha512:      sha512Digest,
						Description: sql.NullString{},
					},
				)
				if err != nil {
					return fmt.Errorf("failed to create file in database: %w", err)
				}
			} else {
				return fmt.Errorf("failed to get file by path %s: %w", filepath.Join(project.Path, file.Name()), err)
			}
		}
		fmt.Printf("Processed file: %s, Size: %d bytes, User: %s\n", fileDb.Name, fileDb.Size, userDb.Username)
	}
	return nil
}

func processingScanDirectory(scanDir string, dbQueries *database.Queries) error {
	wg := &sync.WaitGroup{}
	projects, err := os.ReadDir(scanDir)
	if err != nil {
		return fmt.Errorf("failed to read scan directory: %w", err)
	}
	for _, project := range projects {
		if project.IsDir() {
			wg.Add(1)
			go func() {
				projectPath := filepath.Join(scanDir, project.Name())
				if IsDigitProject(projectPath) {
					userDb, err := GetUserFromPath(projectPath, dbQueries)
					if err != nil {
						log.Fatalf("failed to get user from project path %s: %v", projectPath, err)
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
								log.Fatalf("failed to create project in database: %v", err)
							}
						} else {
							log.Fatalf("failed to get project by path %s: %v", projectPath, err)
						}
					}
					err = processProject(projectDb, dbQueries)
					if err != nil {
						log.Fatalf("failed to process project %s: %v", project.Name(), err)
					}
				}
				wg.Done()
			}()
		}
	}
	wg.Wait()
	return nil
}
