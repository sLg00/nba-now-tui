package client

import (
	"os"
	"testing"
)

func mockPathComponents() PathComponents {
	home, _ := os.UserHomeDir()
	return PathComponents{
		Home:   home,
		Path:   "/foo",
		LLFile: "mock,"}
}

func TestCreateDirectory(t *testing.T) {
	mp := mockPathComponents()
	result, err := createDirectory(mp)
	if result != mp.Home+mp.Path {
		t.Errorf("createDirectory returned %s, expected %s", result, mp)
	}
	if err != nil {
		t.Errorf("createDirectory returned %s, expected no error", err)
	}
	os.Remove(mp.Home + mp.Path)
}

func TestWriteToFiles(t *testing.T) {
	mp := mockPathComponents()
	fileContents := []byte(`{"key": "value"}`)
	result := WriteToFiles(mp.Home+mp.Path+mp.LLFile, fileContents)
	if result != nil {
		t.Errorf("WriteToFiles returned %s, expected no error", result)
	}
	os.Remove(mp.Home + mp.Path + mp.LLFile)
}

func TestFileChecker(t *testing.T) {
	result := fileChecker("/etc/thisfiledoesnotexist")
	if result != false {
		t.Error("Expected false, got ", result)
	}
}
