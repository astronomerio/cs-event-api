package logger

import "log"

type Logger interface {
	Info()
	Debug()
	Error()
	Warn()
}

func NewLogger(kind string) Logger {
	handlers := map[string]func() Logger{
		"mock": func() Logger {
			return NewMockLogger()
		},
	}

	f, ok := handlers[kind]
	if !ok {
		log.Fatal("invalid logger")
	}

	return f()
}
