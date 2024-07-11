package client

import (
	"os"
	"testing"
)

func TestCreateDirectory(t *testing.T) {
	home, _ := os.UserHomeDir()
	mp := PathComponents{
		Home: home,
		Path: "/foo",
	}
	result, err := createDirectory(mp)
	if result != mp.Home+mp.Path {
		t.Errorf("createDirectory returned %s, expected %s", result, mp)
	}
	if err != nil {
		t.Errorf("createDirectory returned %s, expected no error", err)
	}
	os.Remove(mp.Home + mp.Path)
}

func TestFileChecker(t *testing.T) {
	result := fileChecker("/etc/thisfiledoesnotexist")
	if result != false {
		t.Error("Expected false, got ", result)
	}
}
