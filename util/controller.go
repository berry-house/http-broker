package util

import (
	"net/http"

	"go.uber.org/zap"
)

// LogError extracts the Logger from the request's context and logs an error.
func LogError(r *http.Request, err error) {
	loggerValue := r.Context().Value("logger")
	if loggerValue == nil {
		return
	}

	logger, ok := loggerValue.(*zap.Logger)
	if !ok {
		panic("No logger")
	}

	logger.Error(
		err.Error(),
		zap.String("method", r.Method),
		zap.String("url", r.URL.Path),
	)
}

// LogInfo extracts the Logger from the request's context and logs a message.
func LogInfo(r *http.Request, msg string) {
	loggerValue := r.Context().Value("logger")
	if loggerValue == nil {
		return
	}

	logger, ok := loggerValue.(*zap.Logger)
	if !ok {
		panic("No logger")
	}

	logger.Info(
		msg,
		zap.String("method", r.Method),
		zap.String("url", r.URL.Path),
	)
}
