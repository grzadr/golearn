package path

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// PathValidator encapsulates path validation and normalization logic
type PathValidator struct {
	originalPath   string
	normalizedPath string
	isAbsolute     bool
}

// NewPathValidator creates and validates a path
func NewPathValidator(path string) (*PathValidator, error) {
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	// Clean the path to remove any unnecessary separators and dots
	cleanPath := filepath.Clean(path)

	// Convert to absolute path if relative
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	return &PathValidator{
		originalPath:   path,
		normalizedPath: absPath,
		isAbsolute:     filepath.IsAbs(path),
	}, nil
}

func (pv *PathValidator) Validate() error {
	// Check if file exists and is accessible
	fileInfo, err := os.Stat(pv.normalizedPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", pv.normalizedPath)
		}
		return fmt.Errorf("error accessing file: %w", err)
	}

	// Check if it's a regular file
	if !fileInfo.Mode().IsRegular() {
		return fmt.Errorf("path is not a regular file: %s", pv.normalizedPath)
	}

	return nil
}

// GetFS returns an fs.FS for the directory containing the file
func (pv *PathValidator) GetFS() (fs.FS, string, error) {
	// Get the directory containing the file
	dir := filepath.Dir(pv.normalizedPath)

	// Get the base name of the file
	filename := filepath.Base(pv.normalizedPath)

	// Create fs.FS from the directory
	fsys := os.DirFS(dir)

	return fsys, filename, nil
}

func ValidatePath(path *string) (fsys fs.FS, filename string, err error) {
	validator, err := NewPathValidator(*path)
	if err != nil {
		return fsys, filename, fmt.Errorf("Error validating path %s: %w\n", *path, err)
	}

	if err := validator.Validate(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return fsys, filename, fmt.Errorf("Error validating path %s: %w\n", *path, err)
	}

	fsys, filename, err = validator.GetFS()
	if err != nil {
		fmt.Printf("Error creating filesystem: %v\n", err)
		return
	}

	return fsys, filename, err
}
