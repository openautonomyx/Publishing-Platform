package main

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ============================================================================
// Logger Instance (Global)
// ============================================================================

var logger *zap.SugaredLogger

func init() {
	initLogger()
}

// ============================================================================
// Logger Initialization
// ============================================================================

func initLogger() {
	logLevel := os.Getenv("API_LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	var config zap.Config

	if os.Getenv("API_ENV") == "production" {
		// Production configuration
		config = zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(parseLogLevel(logLevel))
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		// Development configuration
		config = zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(parseLogLevel(logLevel))
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	zapLogger, err := config.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	logger = zapLogger.Sugar()
}

func parseLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// ============================================================================
// Helper Functions for Logging
// ============================================================================

func logRequest(method, path string, statusCode int, duration time.Duration, requestID string) {
	logger.Infow("http_request",
		"method", method,
		"path", path,
		"status_code", statusCode,
		"duration_ms", duration.Milliseconds(),
		"request_id", requestID,
	)
}

// ============================================================================
// Structured Logging Helpers
// ============================================================================

func logRateLimitExceeded(clientIP string) {
	logger.Warnw("rate_limit_exceeded",
		"client_ip", clientIP,
	)
}

func logPanic(message string, stack interface{}) {
	logger.Errorw("panic_recovered",
		"message", message,
		"stack", fmt.Sprintf("%v", stack),
	)
}

// ============================================================================
// Time Helper (for consistent time in testing)
// ============================================================================

var timeNow = time.Now

// ============================================================================
// Audit Logging
// ============================================================================
