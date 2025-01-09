package helpers

import "os"

type TestSetup struct {
	originalArgs []string
}

// SetupTest just sets the os.Args to a mocked value set for testing purposes
func SetupTest() *TestSetup {
	ts := &TestSetup{
		originalArgs: os.Args,
	}
	os.Args = []string{"appName", "-d", "2024-12-01"}
	return ts
}

// CleanUpTest reverts os.Args
func (ts *TestSetup) CleanUpTest() {
	os.Args = ts.originalArgs
}
