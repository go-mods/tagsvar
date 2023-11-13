package fs

import (
	"os"
	"path/filepath"
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
