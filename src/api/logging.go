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

	defer zapLogger.Sync()
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

func logError(message string, err error, requestID string) {
	if appErr, ok := IsAppError(err); ok {
		logger.Errorw(message,
			"error_code", appErr.Code,
			"error_message", appErr.Message,
			"status_code", appErr.StatusCode,
			"request_id", requestID,
		)
	} else {
		logger.Errorw(message,
			"error", err.Error(),
			"request_id", requestID,
		)
	}
}

func logDatabaseOperation(operation string, table string, duration time.Duration) {
	logger.Debugw("database_operation",
		"operation", operation,
		"table", table,
		"duration_ms", duration.Milliseconds(),
	)
}

func logServiceCall(service string, method string, duration time.Duration, success bool) {
	logger.Debugw("service_call",
		"service", service,
		"method", method,
		"duration_ms", duration.Milliseconds(),
		"success", success,
	)
}

// ============================================================================
// Structured Logging Helpers
// ============================================================================

func logAuthAttempt(userID, orgID string, success bool) {
	logger.Infow("auth_attempt",
		"user_id", userID,
		"org_id", orgID,
		"success", success,
	)
}

func logDataAccess(userID, orgID string, resourceType string, action string) {
	logger.Infow("data_access",
		"user_id", userID,
		"org_id", orgID,
		"resource_type", resourceType,
		"action", action,
	)
}

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

func logAuditEvent(userID, orgID, action, resourceType string, resourceID interface{}, changes interface{}) {
	logger.Infow("audit_event",
		"user_id", userID,
		"org_id", orgID,
		"action", action,
		"resource_type", resourceType,
		"resource_id", resourceID,
		"changes", changes,
		"timestamp", timeNow().UTC().Format("2006-01-02T15:04:05Z07:00"),
	)
}
