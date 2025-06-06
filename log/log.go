package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Level represents the severity level of a log message
type Level int

const (
	// DEBUG level for detailed troubleshooting information
	DEBUG Level = iota
	// INFO level for general operational information
	INFO
	// WARN level for warning conditions
	WARN
	// ERROR level for error conditions
	ERROR
	// FATAL level for critical errors that cause program termination
	FATAL
)

// String returns the string representation of the log level
func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Color returns the ANSI color code for the log level
func (l Level) Color() string {
	switch l {
	case DEBUG:
		return "\033[36m" // Cyan
	case INFO:
		return "\033[32m" // Green
	case WARN:
		return "\033[33m" // Yellow
	case ERROR:
		return "\033[31m" // Red
	case FATAL:
		return "\033[35m" // Magenta
	default:
		return "\033[0m" // Reset
	}
}

// Subsystem represents a component of the blockchain
type Subsystem string

// Common subsystems in a blockchain
const (
	BLOCKCHAIN Subsystem = "BLOCKCHAIN"
	P2P        Subsystem = "P2P"
	CONSENSUS  Subsystem = "CONSENSUS"
	MEMPOOL    Subsystem = "MEMPOOL"
	RPC        Subsystem = "RPC"
	WALLET     Subsystem = "WALLET"
	DATABASE   Subsystem = "DATABASE"
	MINER      Subsystem = "MINER"
	API        Subsystem = "API"
	GENERAL    Subsystem = "GENERAL"
)

// Logger represents a logger instance
type Logger struct {
	level      Level
	subsystem  Subsystem
	outputs    []io.Writer
	useColors  bool
	showCaller bool
	mu         sync.Mutex
}

// Config represents logger configuration
type Config struct {
	Level      Level
	Subsystem  Subsystem
	UseColors  bool
	ShowCaller bool
	LogToFile  bool
	LogDir     string
	FileName   string
}

// DefaultConfig returns a default logger configuration
func DefaultConfig() Config {
	return Config{
		Level:      INFO,
		Subsystem:  GENERAL,
		UseColors:  true,
		ShowCaller: true,
		LogToFile:  false,
		LogDir:     "logs",
		FileName:   "",
	}
}

// New creates a new logger with the given configuration
func New(config Config) (*Logger, error) {
	outputs := []io.Writer{os.Stdout}

	if config.LogToFile {
		if err := os.MkdirAll(config.LogDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		fileName := config.FileName
		if fileName == "" {
			fileName = fmt.Sprintf("%s_%s.log",
				strings.ToLower(string(config.Subsystem)),
				time.Now().Format("2006-01-02"))
		}

		filePath := filepath.Join(config.LogDir, fileName)
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}

		outputs = append(outputs, file)
	}

	return &Logger{
		level:      config.Level,
		subsystem:  config.Subsystem,
		outputs:    outputs,
		useColors:  config.UseColors,
		showCaller: config.ShowCaller,
		mu:         sync.Mutex{},
	}, nil
}

// NewWithSubsystem creates a new logger for a specific subsystem
func NewWithSubsystem(subsystem Subsystem) (*Logger, error) {
	config := DefaultConfig()
	config.Subsystem = subsystem
	return New(config)
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// AddOutput adds an additional output writer
func (l *Logger) AddOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.outputs = append(l.outputs, w)
}

// log logs a message with the specified level
func (l *Logger) log(level Level, format string, args ...any) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	message := fmt.Sprintf(format, args...)

	var caller string
	if l.showCaller {
		_, file, line, ok := runtime.Caller(2)
		if ok {
			file = filepath.Base(file)
			caller = fmt.Sprintf(" [%s:%d]", file, line)
		}
	}

	var logLine string
	if l.useColors {
		reset := "\033[0m"
		logLine = fmt.Sprintf("%s %s[%s]%s [%s]%s%s %s\n",
			timestamp,
			level.Color(), level, reset,
			l.subsystem,
			caller,
			reset,
			message)
	} else {
		logLine = fmt.Sprintf("%s [%s] [%s]%s %s\n",
			timestamp,
			level,
			l.subsystem,
			caller,
			message)
	}

	for _, output := range l.outputs {
		fmt.Fprint(output, logLine)
	}

	// If level is FATAL, exit the program
	if level == FATAL {
		os.Exit(1)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...any) {
	l.log(DEBUG, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...any) {
	l.log(INFO, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...any) {
	l.log(WARN, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...any) {
	l.log(ERROR, format, args...)
}

// Fatal logs a fatal message and exits the program
func (l *Logger) Fatal(format string, args ...any) {
	l.log(FATAL, format, args...)
}

// Global logger instances for each subsystem
var (
	globalMu      sync.Mutex
	globalLoggers = make(map[Subsystem]*Logger)
)

// GetLogger returns a logger for the specified subsystem
func GetLogger(subsystem Subsystem) *Logger {
	globalMu.Lock()
	defer globalMu.Unlock()

	if logger, ok := globalLoggers[subsystem]; ok {
		return logger
	}

	logger, err := NewWithSubsystem(subsystem)
	if err != nil {
		// Fallback to a basic logger if there's an error
		logger = &Logger{
			level:      INFO,
			subsystem:  subsystem,
			outputs:    []io.Writer{os.Stdout},
			useColors:  true,
			showCaller: true,
			mu:         sync.Mutex{},
		}
	}

	globalLoggers[subsystem] = logger
	return logger
}

// ConfigureGlobalLoggers configures all global loggers with the same settings
func ConfigureGlobalLoggers(level Level, logToFile bool, logDir string) error {
	globalMu.Lock()
	defer globalMu.Unlock()

	for subsystem, logger := range globalLoggers {
		logger.SetLevel(level)

		if logToFile {
			if err := os.MkdirAll(logDir, 0755); err != nil {
				return fmt.Errorf("failed to create log directory: %w", err)
			}

			fileName := fmt.Sprintf("%s_%s.log",
				strings.ToLower(string(subsystem)),
				time.Now().Format("2006-01-02"))

			filePath := filepath.Join(logDir, fileName)
			file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return fmt.Errorf("failed to open log file: %w", err)
			}

			logger.AddOutput(file)
		}
	}

	return nil
}
