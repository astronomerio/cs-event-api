package logger

type MockLogger struct{}

func NewMockLogger() *MockLogger {
	return &MockLogger{}
}

func (l *MockLogger) Info() {
	// noop
}

func (l *MockLogger) Debug() {
	// noop
}

func (l *MockLogger) Error() {
	// noop
}

func (l *MockLogger) Warn() {
	// noop
}
