package logger

import (
	"os"
	"sync"
)

type Logger struct {
	file *os.File
	mu   sync.Mutex
}

func NewLogger(filename string) (*Logger, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return &Logger{file: f}, nil
}

func (l *Logger) Log(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.file.WriteString(msg + "\n")
}

func (l *Logger) Close() {
	l.file.Close()
}
