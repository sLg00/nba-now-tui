package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLogToFile(t *testing.T) {
	home, _ := os.UserHomeDir()
	logDir := filepath.Join(home, ".config/nba-tui/logs/")
	fileName := filepath.Join(logDir, "appLog.log")
	result, _ := LogToFile()
	if fileName != result {
		t.Errorf("LogToFile() failed: expected %s, got %s", fileName, result)
	}

}
