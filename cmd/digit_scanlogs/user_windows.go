//go:build windows

package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"golang.org/x/sys/windows"

	"github.com/edipretoro/digit_scanlogs/internal/database"
	"github.com/google/uuid"
)

func GetUserFromPath(path string, dbQueries *database.Queries) (userDb database.User, err error) {
	f, err := os.Open(path)
	if err != nil {
		return userDb, err
	}
	defer f.Close()

	handle := windows.Handle(f.Fd())
	var sd *windows.SECURITY_DESCRIPTOR
	sd, err = windows.GetSecurityInfo(
		handle,
		windows.SE_FILE_OBJECT,
		windows.OWNER_SECURITY_INFORMATION,
	)
	if err != nil {
		return userDb, err
	}

	var owner *windows.SID
	owner, _, err = sd.Owner()
	if err != nil {
		return userDb, err
	}

	uid := int64(owner.SubAuthority(0))

	account, domain, _, err := owner.LookupAccount("")
	if err != nil {
		return userDb, err
	}
	userDb, err = dbQueries.GetUserByUsername(context.Background(), account)
	if err != nil {
		if err == sql.ErrNoRows {
			userDb, err = dbQueries.CreateUser(
				context.Background(),
				database.CreateUserParams{
					ID:       uuid.NewString(),
					Uid:      uid,
					Username: account,
					Fullname: fmt.Sprintf("%s\\%s", domain, account),
				},
			)
			if err != nil {
				return database.User{}, fmt.Errorf("failed to create user in database: %w", err)
			}
		} else {
			return database.User{}, fmt.Errorf("failed to get user by username %s (%d): %w", owner.String(), uid, err)
		}
	}
	return userDb, nil
}
