package internal

import (
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func createMockDir() (string, error) {
	tmpDir, err := os.MkdirTemp("", "test_findfiles")
	if err != nil {
		log.Println(err)
		return "", err
	}
	return tmpDir, nil
}

func createMockFiles() ([]string, string, error) {
	today := time.Now().Format("2006-01-02")
	files := []struct {
		name    string
		content string
	}{
		{today + "_file", "abc"},
		{"2024-01-02_file", "xyz"},
		{"2024-01-03_file", "qwe"},
	}
	var fileResults []string
	dir, _ := createMockDir()
	for _, file := range files {
		filePath := filepath.Join(dir, file.name)
		if err := os.WriteFile(filePath, []byte(file.content), 0644); err != nil {
			log.Printf("error creating file %s: %v", file.name, err)
			return nil, "", err
		}
		fileResults = append(fileResults, filePath)
	}
	return fileResults, dir, nil
}

func TestFindFiles(t *testing.T) {
	_, mockDir, err := createMockFiles()
	if err != nil {
		t.Errorf("error creating mock files: %v", err)
	}
	defer os.RemoveAll(mockDir)

	pattern := `^(\d{4}-\d{2}-\d{2})_file$`

	results, err := FindFiles(mockDir, pattern)
	if err != nil {
		t.Fatalf("Find files returned an error: %v", err)
	}

	expectedFiles := []string{
		filepath.Join(mockDir, "2024-01-02_file"),
		filepath.Join(mockDir, "2024-01-03_file"),
	}

	if len(results) != len(expectedFiles) {
		t.Errorf("Find files returned %d files, expected %d, files returned %s", len(results), len(expectedFiles), results)
	}

	expectedSet := make(map[string]bool)
	for _, file := range expectedFiles {
		expectedSet[file] = true
	}
	for _, result := range results {
		if !expectedSet[result] {
			t.Errorf("unexpected file found: %s", result)
		}
	}
}

func TestRemoveFiles(t *testing.T) {
	//

}
