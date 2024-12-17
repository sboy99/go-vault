package utils

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ReadFile reads the content of a file and returns it as a string.
func ReadFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// WriteFile writes content to a file, creating the file if it does not exist.
func WriteFile(filePath string, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

// AppendToFile appends content to a file, creating the file if it does not exist.
func AppendToFile(filePath string, content string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

// DeleteFile deletes the specified file.
func DeleteFile(filePath string) error {
	return os.Remove(filePath)
}

// FileExists checks if a file exists at the given path.
func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// DirectoryExists checks if a directory exists at the given path.
func DirectoryExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// CreateDirectory creates a new directory, including any necessary parents.
func CreateDirectory(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
}

// CreateDirIfNotExist creates a new directory if it does not already exist.
func CreateDirIfNotExist(dirPath string) error {
	if IsPathContainsFile(dirPath) {
		dirPath = filepath.Dir(dirPath)
	}
	if !DirectoryExists(dirPath) {
		return CreateDirectory(dirPath)
	}
	return nil
}

// ListFiles lists all files in a given directory.
func ListFiles(dirPath string) ([]string, error) {
	files := []string{}
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

// GetFileSize returns the size of the file in bytes.
func GetFileSize(filePath string) (int64, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// GetFileExtension returns the file extension of the given file.
func GetFileExtension(filePath string) string {
	return filepath.Ext(filePath)
}

// CopyFile copies a file from the source path to the destination path.
func CopyFile(srcPath, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}

// RenameFile renames a file from oldPath to newPath.
func RenameFile(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

// IsPathContainsFile checks if the path contains a file.
func IsPathContainsFile(path string) bool {
	dirs := strings.Split(path, "/")
	last := dirs[len(dirs)-1]
	return strings.Contains(last, ".ts")
}
