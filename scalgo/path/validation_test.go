package path

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewPathValidator(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		expectError   bool
		errorContains string
	}{
		{
			name:          "Empty path",
			path:          "",
			expectError:   true,
			errorContains: "path cannot be empty",
		},
		{
			name:        "Valid relative path",
			path:        "testdata/test.txt",
			expectError: false,
		},
		{
			name:        "Valid absolute path",
			path:        filepath.Join(os.TempDir(), "test.txt"),
			expectError: false,
		},
		{
			name:        "Path with unnecessary separators",
			path:        "testdata///test.txt",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator, err := NewPathValidator(tt.path)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error containing '%s', got nil", tt.errorContains)
					return
				}
				if tt.errorContains != "" && !errors.Is(err, err) {
					t.Errorf("Expected error containing '%s', got '%v'", tt.errorContains, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if validator == nil {
				t.Error("Expected validator to not be nil")
				return
			}

			if validator.originalPath != tt.path {
				t.Errorf("Expected originalPath to be '%s', got '%s'", tt.path, validator.originalPath)
			}
		})
	}
}

func TestPathValidator_Validate(t *testing.T) {
	// Create temporary test files
	tempDir := t.TempDir()
	regularFile := filepath.Join(tempDir, "regular.txt")
	if err := os.WriteFile(regularFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a directory for testing non-regular file case
	dirPath := filepath.Join(tempDir, "testdir")
	if err := os.Mkdir(dirPath, 0755); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name          string
		path          string
		expectError   bool
		errorContains string
		setup         func() string
		cleanup       func(string) error
	}{
		{
			name:        "Valid regular file",
			path:        regularFile,
			expectError: false,
		},
		{
			name:          "Non-existent file",
			path:          filepath.Join(tempDir, "nonexistent.txt"),
			expectError:   true,
			errorContains: "file does not exist",
		},
		{
			name:          "Directory instead of file",
			path:          dirPath,
			expectError:   true,
			errorContains: "path is not a regular file",
		},
		{
			name:          "Permission denied",
			expectError:   true,
			errorContains: "error accessing file",
			setup: func() string {
				// Create a file in a subdirectory with no read permissions
				noPermDir := filepath.Join(tempDir, "noperm")
				if err := os.Mkdir(noPermDir, 0755); err != nil {
					t.Fatal(err)
				}

				noReadFile := filepath.Join(noPermDir, "noperm.txt")
				if err := os.WriteFile(noReadFile, []byte("test"), 0644); err != nil {
					t.Fatal(err)
				}

				// First remove all permissions from the file
				if err := os.Chmod(noReadFile, 0000); err != nil {
					t.Fatal(err)
				}

				// Then remove permissions from the directory
				if err := os.Chmod(noPermDir, 0000); err != nil {
					t.Fatal(err)
				}

				return noReadFile
			},
			cleanup: func(path string) error {
				// Restore directory permissions first
				dir := filepath.Dir(path)
				if err := os.Chmod(dir, 0755); err != nil {
					return fmt.Errorf("failed to restore directory permissions: %w", err)
				}

				// Then restore file permissions
				if err := os.Chmod(path, 0644); err != nil {
					return fmt.Errorf("failed to restore file permissions: %w", err)
				}

				// Clean up file
				if err := os.Remove(path); err != nil {
					return fmt.Errorf("failed to remove test file: %w", err)
				}

				// Clean up directory
				if err := os.Remove(dir); err != nil {
					return fmt.Errorf("failed to remove test directory: %w", err)
				}

				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testPath string
			if tt.setup != nil {
				testPath = tt.setup()
			} else {
				testPath = tt.path
			}

			// Ensure cleanup runs even if the test fails
			defer func() {
				if tt.cleanup != nil {
					if err := tt.cleanup(testPath); err != nil {
						t.Logf("Cleanup failed: %v", err)
					}
				}
			}()

			validator, err := NewPathValidator(testPath)
			if err != nil {
				t.Fatalf("Failed to create validator: %v", err)
			}

			err = validator.Validate()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error containing '%s', got nil", tt.errorContains)
					return
				}
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error containing '%s', got '%v'", tt.errorContains, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestPathValidator_GetFS(t *testing.T) {
	// Create temporary test files
	tempDir := t.TempDir()
	regularFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(regularFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		path           string
		expectedFile   string
		expectError    bool
		validateFsFunc func(fs.FS, string) error
	}{
		{
			name:         "Valid file path",
			path:         regularFile,
			expectedFile: "test.txt",
			validateFsFunc: func(fsys fs.FS, filename string) error {
				// Try to open the file using the returned FS
				f, err := fsys.Open(filename)
				if err != nil {
					return err
				}
				defer f.Close()
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator, err := NewPathValidator(tt.path)
			if err != nil {
				t.Fatalf("Failed to create validator: %v", err)
			}

			fsys, filename, err := validator.GetFS()
			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if filename != tt.expectedFile {
				t.Errorf("Expected filename '%s', got '%s'", tt.expectedFile, filename)
			}

			if tt.validateFsFunc != nil {
				if err := tt.validateFsFunc(fsys, filename); err != nil {
					t.Errorf("FS validation failed: %v", err)
				}
			}
		})
	}
}

func TestValidatePath(t *testing.T) {
	// Create temporary test files
	tempDir := t.TempDir()
	regularFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(regularFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name          string
		path          string
		expectError   bool
		errorContains string
	}{
		{
			name:        "Valid path",
			path:        regularFile,
			expectError: false,
		},
		{
			name:          "Empty path",
			path:          "",
			expectError:   true,
			errorContains: "path cannot be empty",
		},
		{
			name:          "Non-existent path",
			path:          filepath.Join(tempDir, "nonexistent.txt"),
			expectError:   true,
			errorContains: "file does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.path
			fsys, filename, err := ValidatePath(&path)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error containing '%s', got nil", tt.errorContains)
					return
				}
				if tt.errorContains != "" && !errors.Is(err, err) {
					t.Errorf("Expected error containing '%s', got '%v'", tt.errorContains, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if fsys == nil {
				t.Error("Expected non-nil filesystem")
			}

			if filename == "" {
				t.Error("Expected non-empty filename")
			}

			// Verify the returned fs.FS works
			f, err := fsys.Open(filename)
			if err != nil {
				t.Errorf("Failed to open file using returned FS: %v", err)
				return
			}
			defer f.Close()
		})
	}
}

func TestNewPathValidator_AbsolutePathError(t *testing.T) {
	// Save current directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Create a directory and remove it to cause Abs to fail
	tempDir := t.TempDir()
	workDir := filepath.Join(tempDir, "nonexistent")
	if err := os.Mkdir(workDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Change to the directory and then remove it
	if err := os.Chdir(workDir); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(workDir); err != nil {
		t.Fatal(err)
	}

	// Try to create a PathValidator - should fail because working directory doesn't exist
	_, err = NewPathValidator("test.txt")

	// Restore working directory
	if err := os.Chdir(originalWd); err != nil {
		t.Fatal(err)
	}

	if err == nil {
		t.Error("Expected error about absolute path resolution, got nil")
	}
	if !strings.Contains(err.Error(), "failed to resolve absolute path") {
		t.Errorf("Expected error about absolute path resolution, got: %v", err)
	}
}
