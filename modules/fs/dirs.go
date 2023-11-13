package fs

import (
	"os"
	"path/filepath"
)

// WorkDir returns the working directory
func WorkDir(root ...string) (string, error) {
	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// If the root is not specified, the current directory is used
	if len(root) == 0 || root[0] == "" || root[0] == "." {
		return workDir(cwd, cwd)
	}

	return workDir(cwd, root[0])
}

// workDir returns the working directory
func workDir(cwd string, path string) (string, error) {
	// If the path has no volume name, the current directory is used
	// and the path is relative to the current directory
	if filepath.VolumeName(path) == "" {
		path = filepath.Join(cwd, path)
	}

	// Convert the path to absolute
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	// If the path is a directory, it is used as the working directory
	isDir, err := IsDir(path)
	if err != nil {
		return "", err
	}
	if isDir {
		return path, nil
	}

	// If the path is a file, the directory of the path is used as the working directory
	return filepath.Dir(path), nil
}

// IsDir returns true if the path is a directory
func IsDir(path string) (bool, error) {
	// Get the file info
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	// Return true if the path is a directory
	return info.IsDir(), nil
}
