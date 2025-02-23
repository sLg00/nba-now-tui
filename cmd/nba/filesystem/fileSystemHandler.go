package filesystemops

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileSystemHandler interface {
	WriteFile(file string, data []byte) error
	ReadFile(file string) ([]byte, error)
	FileExists(file string) bool
	CleanOldFiles(pc []string) error
}

type DefaultFsHandler struct {
	baseDirectory string
}

func NewDefaultFsHandler() *DefaultFsHandler {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Errorf("could not get home directory: %v", err)
		return nil
	}
	return &DefaultFsHandler{baseDirectory: home}
}

func (fs *DefaultFsHandler) WriteFile(file string, data []byte) error {
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory failed: %w", err)
	}

	return os.WriteFile(file, data, 0644)
}

func (fs *DefaultFsHandler) ReadFile(file string) ([]byte, error) {
	workFile, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("reading file failed: %w", err)
	}
	return workFile, nil
}

func (fs *DefaultFsHandler) FileExists(file string) bool {
	fileInfo, err := os.Stat(file)
	if err == nil {
		if fileInfo.Size() > 1000 {
			return true
		}
	}
	return false
}

func (fs *DefaultFsHandler) CleanOldFiles(pc []string) error {

	filesRegex := "^(\\d{4}-\\d{2}-\\d{2})_.*$"

	for _, path := range pc {
		fileList, err := FindFiles(path, filesRegex)
		if err != nil {
			return fmt.Errorf("could not list files in path %s: %v", path, err)
		}
		err = RemoveFiles(fileList)
		if err != nil {
			return fmt.Errorf("could not remove files in path %s: %v", path, err)
		}
	}

	return nil
}
