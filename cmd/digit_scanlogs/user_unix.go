//go:build !windows

package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"

	"github.com/edipretoro/digit_scanlogs/internal/database"
	"github.com/google/uuid"
)

func GetUserFromPath(path string, dbQueries *database.Queries) (userDb database.User, err error) {
	metadata, err := os.Stat(path)
	if err != nil {
		return database.User{}, fmt.Errorf("failed to get metadata for path %s: %w", path, err)
	}
	metadataStat := metadata.Sys().(*syscall.Stat_t)
	userUid := strconv.FormatUint(uint64(metadataStat.Uid), 10)
	user, err := user.LookupId(userUid)
	if err != nil {
		return database.User{}, fmt.Errorf("failed to lookup user for UID %d: %w", metadataStat.Uid, err)
	}
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
				return database.User{}, fmt.Errorf("failed to create user in database: %w", err)
			}
		} else {
			return database.User{}, fmt.Errorf("failed to get user by UID %d: %w", metadataStat.Uid, err)
		}
	}
	return userDb, nil
}
