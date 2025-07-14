package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"github.com/edipretoro/digit_scanlogs/internal/database"
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
				metadata, err := os.Stat(projectPath)
				if err != nil {
					return fmt.Errorf("failed to get metadata for project %s: %w", project.Name(), err)
				}
				metadataStat := metadata.Sys().(*syscall.Stat_t)
				fmt.Printf("Found Digit project: %s\n", project.Name())
				fmt.Printf("Project metadata: UID=%d, GID=%d, Size=%d bytes\n", metadataStat.Uid, metadataStat.Gid, metadata.Size())
				userUid := strconv.FormatUint(uint64(metadataStat.Uid), 10)
				user, err := user.LookupId(userUid)
				if err != nil {
					return fmt.Errorf("failed to lookup user for UID %d: %w", metadataStat.Uid, err)
				}
				var userDb database.User
				userDb, err = dbQueries.GetUserByUID(context.Background(), int64(metadataStat.Uid))
				if err != nil {
					if err == sql.ErrNoRows {
						userDb, err = dbQueries.CreateUser(
							context.Background(),
							database.CreateUserParams{
								ID:       uuid.NewString(),
								Uid:      int64(metadataStat.Uid),
								Username: user.Username,
								Fullname: user.Name,
							},
						)
						if err != nil {
							return fmt.Errorf("failed to create user in database: %w", err)
						}
					} else {
						return fmt.Errorf("failed to get user by UID %d: %w", metadataStat.Uid, err)
					}
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
				err = processProject(projectDb)
				if err != nil {
					return fmt.Errorf("failed to process project %s: %w", project.Name(), err)
				}
			}
		}
	}
	return nil
}
