package fs

import (
	"github.com/go-mods/tagsvar/modules/config"
	"os"
	"path/filepath"
	"strings"
)

// ListFiles lists files in a directory
// The file name must be checked through the filter function
// If recursive is true, the files are listed in all subdirectories
func ListFiles(dir string, recursive bool, filter func(string) bool) ([]string, error) {
	files := make([]string, 0)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories if recursive is false
		if info.IsDir() && !recursive {
			return filepath.SkipDir
		}

		// Get the file name
		fileName := info.Name()

		// Skip files that do not match the filter function
		if !filter(fileName) {
			return nil
		}

		files = append(files, path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// IsGoFile checks if the file is a .go file
func IsGoFile(fileName string) bool {
	return filepath.Ext(fileName) == ".go"
}

// IsGeneratedFile checks if the file is a generated file
// It must be a .go file
// It must start with the prefix and end with the suffix
func IsGeneratedFile(fileName string) bool {
	// Get the file name only
	fileName = filepath.Base(fileName)
	// Check if the file is a .go file
	if !IsGoFile(fileName) {
		return false
	}
	// Remove the extension
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))
	// Check if the file starts with the prefix
	if !strings.HasPrefix(filepath.Base(fileName), config.C.Prefix) {
		return false
	}
	// Check if the file ends with the suffix
	if !strings.HasSuffix(filepath.Base(fileName), config.C.Suffix) {
		return false
	}
	return true
}

// IsProjectFile checks if the file is a project file (not a generated file, neither a test file)
func IsProjectFile(fileName string) bool {
	// Get the file name only
	fileName = filepath.Base(fileName)
	// Check if the file is a .go file
	if !IsGoFile(fileName) {
		return false
	}
	// Check if the file is a generated file
	if IsGeneratedFile(fileName) {
		return false
	}
	// Check if the file is a test file
	if strings.HasSuffix(fileName, "_test.go") {
		return false
	}
	return true
}
