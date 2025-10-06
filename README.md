# digit_scanlogs

A Go application for tracking and monitoring digitized document projects. This tool scans directories containing digitizing projects, catalogs files, computes SHA-512 checksums for integrity verification, and maintains a SQLite database of all scanned documents.

## Overview

`digit_scanlogs` helps manage large collections of digitized documents by maintaining a comprehensive database of files, their metadata, and cryptographic checksums. The tool automatically tracks file ownership, project organization, and file integrity across Unix-like systems and Windows.

### Key Features

- **Automatic Project Detection**: Identifies digitizing projects by looking for `.mets` files
- **File Integrity Tracking**: Computes and stores SHA-512 checksums for all files
- **User Tracking**: Associates files with their system owners (UID-based on Unix, alternative on Windows)
- **Concurrent Processing**: Uses goroutines for efficient parallel file scanning
- **Cross-Platform Support**: Works on Linux, macOS (darwin), and Windows
- **SQLite Database**: Lightweight, embedded database for metadata storage

## Prerequisites

- **Go**: Version 1.24.0 or later
- **Goose**: Database migration tool (optional, for schema management)
- **SQLite**: Built-in via modernc.org/sqlite (no external dependencies)

## Installation

### From Source

Clone the repository and build the application:

```bash
git clone https://github.com/edipretoro/digit_scanlogs.git
cd digit_scanlogs
go mod download
go build -o digitcheck ./cmd/digit_check
```

### Cross-Platform Builds

The project includes pre-built binaries in the `build/` directory for multiple platforms:

- `digitcheck-darwin-arm64` - macOS (Apple Silicon)
- `digitcheck-linux-amd64` - Linux (64-bit)
- `digitcheck-windows-amd64.exe` - Windows (64-bit)

To create these builds manually:

```bash
# macOS (ARM64)
GOOS=darwin GOARCH=arm64 go build -o build/digitcheck-darwin-arm64 ./cmd/digit_check

# Linux (AMD64)
GOOS=linux GOARCH=amd64 go build -o build/digitcheck-linux-amd64 ./cmd/digit_check

# Windows (AMD64)
GOOS=windows GOARCH=amd64 go build -o build/digitcheck-windows-amd64.exe ./cmd/digit_check
```

## Configuration

The application requires a `.env` file in the project root with the following variables:

```env
# Database connection string for the application
DIGIT_SCAN_DSN=/path/to/digit_scanlogs.db?_pragma=busy_timeout%3D10000&_pragma=journal_mode%3DWAL&mode=rwc

# Directory containing DIGIT projects to scan
DIGIT_SCAN_DIR=/path/to/your/documents

# Goose migration settings (if using migrations)
GOOSE_DRIVER=sqlite3
GOOSE_DBSTRING=./sql/db/digit_scanlogs.db?_pragma=busy_timeout%3D10000&_pragma=journal_mode%3DWAL&mode=rwc
GOOSE_MIGRATION_DIR=sql/schema
```

### Configuration Parameters

- **DIGIT_SCAN_DSN**: SQLite database connection string with Write-Ahead Logging (WAL) mode enabled
- **DIGIT_SCAN_DIR**: Root directory to scan for DIGIT projects
- **GOOSE_DRIVER**: Database driver for migrations (sqlite3)
- **GOOSE_DBSTRING**: Database connection string for Goose migrations
- **GOOSE_MIGRATION_DIR**: Directory containing migration files

### Example Configuration

Create a `.env` file:

```bash
cat > .env << 'EOF'
GOOSE_DRIVER=sqlite3
GOOSE_DBSTRING=./sql/db/digit_scanlogs.db?_pragma=busy_timeout%3D10000&_pragma=journal_mode%3DWAL&mode=rwc
GOOSE_MIGRATION_DIR=sql/schema
DIGIT_SCAN_DSN=./sql/db/digit_scanlogs.db?_pragma=busy_timeout%3D10000&_pragma=journal_mode%3DWAL&mode=rwc
DIGIT_SCAN_DIR=/path/to/your/scanned/documents
EOF
```

## Database Setup

The application uses SQLite with the following schema:

### Tables

**users**: System users who own files
- `id` (UUID, primary key)
- `uid` (integer, system user ID)
- `username` (text, unique)
- `fullname` (text)
- Timestamps: `created_at`, `updated_at`, `deleted_at`

**projects**: Digitizing project directories
- `id` (UUID, primary key)
- `name` (text, project name)
- `path` (text, filesystem path)
- `description` (text, optional)
- `created_by` (UUID, references users)
- Timestamps: `created_at`, `updated_at`, `deleted_at`

**files**: Individual files within projects
- `id` (UUID, primary key)
- `project_id` (UUID, references projects)
- `user_id` (UUID, references users)
- `name` (text, filename)
- `path` (text, full filesystem path)
- `size` (bigint, file size in bytes)
- `mode` (text, file permissions)
- `modtime` (timestamp, modification time)
- `sha512` (text, SHA-512 checksum)
- `description` (text, optional)
- Timestamps: `created_at`, `updated_at`, `deleted_at`

### Database Initialization

The database is automatically created on first run. If you're using Goose for migrations and have migration files in `sql/schema/`, run:

```bash
goose up
```

## Usage

### Running the Scanner

After configuring the `.env` file, run the scanner:

```bash
./digitcheck
```

Or if built from source:

```bash
go run ./cmd/digit_check
```

### What It Does

1. **Loads Configuration**: Reads environment variables from `.env`
2. **Connects to Database**: Opens or creates the SQLite database
3. **Scans Directory**: Recursively walks through `DIGIT_SCAN_DIR`
4. **Identifies DIGIT Projects**: Looks for directories containing `{project_name}.mets` files
5. **Catalogs Files**: For each file in a DIGIT project:
   - Retrieves file metadata (size, permissions, modification time)
   - Computes SHA-512 checksum
   - Identifies file owner from filesystem
   - Creates or updates database records
6. **Reports Progress**: Logs new files found during the scan

### Example Output

```
2024/07/25 15:04:23 Checking directory: /Users/edipretoro/Documents/maps
2024/07/25 15:04:24 File /Users/edipretoro/Documents/maps/project1/file1.tif created successfully in database
2024/07/25 15:04:25 File /Users/edipretoro/Documents/maps/project1/file2.tif created successfully in database
2024/07/25 15:04:26 Scan directory checked successfully: 2 new files found
```

## Architecture

### Project Structure

```
digit_scanlogs/
├── cmd/
│   └── digit_check/          # Main application
│       ├── main.go            # Entry point
│       └── check.go           # Scanning logic
├── internal/
│   ├── digestfile/            # SHA-512 hashing utilities
│   │   └── digest.go
│   └── user/                  # User management (platform-specific)
│       ├── user_unix.go       # Unix/Linux/macOS implementation
│       └── user_windows.go    # Windows implementation
├── sql/
│   └── db/                    # SQLite database file
│       └── digit_scanlogs.db
├── build/                     # Compiled binaries
├── .env                       # Configuration (not in git)
├── .gitignore
├── go.mod
└── README.md
```

### Key Components

**Main Scanner** (`cmd/digit_check/main.go`)
- Initializes database connection
- Orchestrates the scanning process
- Manages goroutine synchronization

**Checking Logic** (`cmd/digit_check/check.go`)
- Walks directory tree
- Identifies DIGIT projects
- Manages concurrent file processing
- Creates database records

**Digest Calculation** (`internal/digestfile/digest.go`)
- Computes SHA-512 checksums for files
- Used for integrity verification

**User Management** (`internal/user/`)
- Platform-specific user lookup
- Maps filesystem ownership to database records
- Unix: Uses UID from file metadata
- Windows: Alternative implementation

## Digitizing Project Format

A directory is recognized as a digitizing project if it contains a `.mets` file matching the directory name:

```
/path/to/projects/
└── my_project/
    ├── my_project.mets    # This file marks it as a digitizing project
    ├── file1.tif
    ├── file2.tif
    └── metadata.xml
```

## Development

### Building

```bash
# Build for current platform
go build -o digitcheck ./cmd/digit_check

# Run without building
go run ./cmd/digit_check

# Build with all dependencies
go build -v ./...
```

### Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Dependencies

The project uses minimal external dependencies:

- `modernc.org/sqlite` - Pure Go SQLite implementation
- `github.com/joho/godotenv` - Environment variable loading
- `github.com/google/uuid` - UUID generation

Install dependencies:

```bash
go mod download
```

## Performance Considerations

- **Concurrent Processing**: Uses goroutines for parallel file processing
- **SQLite Optimizations**: 
  - Write-Ahead Logging (WAL) mode for better concurrency
  - Busy timeout of 10 seconds to handle concurrent access
  - Single connection limit to prevent database locks
- **Incremental Scanning**: Only processes files not already in the database

## Troubleshooting

### Database Locked Errors

If you encounter "database is locked" errors:
- Ensure only one instance is running
- Check that WAL mode is enabled in the connection string
- Verify the `busy_timeout` pragma is set

### Permission Errors

On Unix systems, ensure the application has read access to:
- All files in the scan directory
- The database file and its parent directory

### Missing Files in Database

Files will only be added if they are:
- Located within a valid digitizing project directory
- Readable by the scanning user
- Not already present in the database

## License

This project is maintained by edipretoro. For licensing information, please contact the repository owner.

## Contributing

Contributions are welcome! Please ensure:
- Code follows Go best practices and formatting (`go fmt`)
- Tests pass (`go test ./...`)
- Commit messages are descriptive

## Contact

For questions or issues, please open an issue on the GitHub repository.
