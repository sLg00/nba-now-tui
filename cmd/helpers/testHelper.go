package helpers

import "os"

type testSetup struct {
	originalArgs []string
}

// SetupTest just sets the os.Args to a mocked value set for testing purposes
func SetupTest() *testSetup {
	ts := &testSetup{
		originalArgs: os.Args,
	}
	os.Args = []string{"appName", "-d", "2024-12-01"}
	return ts
}

// CleanUpTest reverts os.Args
func (ts *testSetup) CleanUpTest() {
	os.Args = ts.originalArgs
}
