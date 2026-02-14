package helpers

type TestSetup struct{}

func SetupTest() *TestSetup {
	return &TestSetup{}
}

func (ts *TestSetup) CleanUpTest() {}
