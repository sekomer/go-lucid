package log_test

import (
	"go-lucid/log"
	"os"
	"path/filepath"
	"testing"
)

func TestLogger(t *testing.T) {
	t.Parallel()

	// Test basic logger functionality
	t.Run("BasicLogging", func(t *testing.T) {
		config := log.DefaultConfig()
		config.Level = log.DEBUG
		config.Subsystem = log.BLOCKCHAIN

		logger, err := log.New(config)
		if err != nil {
			t.Fatalf("Failed to create logger: %v", err)
		}

		// Log messages at different levels
		logger.Debug("This is a debug message: %d", 1)
		logger.Info("This is an info message: %s", "test")
		logger.Warn("This is a warning message")
		logger.Error("This is an error message: %v", map[string]string{"key": "value"})
		// logger.Fatal("This is a fatal message") // This would exit the program
	})

	// Test file logging
	t.Run("FileLogging", func(t *testing.T) {
		// Create a temporary directory for logs
		tempDir := filepath.Join(os.TempDir(), "lucid-logs-test")
		defer os.RemoveAll(tempDir)

		config := log.DefaultConfig()
		config.Level = log.INFO
		config.Subsystem = log.P2P
		config.LogToFile = true
		config.LogDir = tempDir
		config.FileName = "p2p_test.log"

		logger, err := log.New(config)
		if err != nil {
			t.Fatalf("Failed to create logger: %v", err)
		}

		logger.Info("This message should go to both console and file")

		// Verify the log file exists
		logFile := filepath.Join(tempDir, "p2p_test.log")
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			t.Errorf("Log file was not created: %v", err)
		}
	})

	// Test global loggers
	t.Run("GlobalLoggers", func(t *testing.T) {
		// Get loggers for different subsystems
		blockchainLogger := log.GetLogger(log.BLOCKCHAIN)
		p2pLogger := log.GetLogger(log.P2P)
		consensusLogger := log.GetLogger(log.CONSENSUS)

		blockchainLogger.Info("Blockchain info message")
		p2pLogger.Info("P2P info message")
		consensusLogger.Info("Consensus info message")

		// Configure all loggers to use DEBUG level
		err := log.ConfigureGlobalLoggers(log.DEBUG, false, "")
		if err != nil {
			t.Fatalf("Failed to configure global loggers: %v", err)
		}

		blockchainLogger.Debug("This debug message should now be visible")
	})
}
