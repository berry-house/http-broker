package util

import (
	"net/http"

	"go.uber.org/zap"
)

func LogError(r *http.Request, err error) {
	loggerValue := r.Context().Value("logger")
	if loggerValue == nil {
		panic("No logger")
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

func LogInfo(r *http.Request, msg string) {
	loggerValue := r.Context().Value("logger")
	if loggerValue == nil {
		panic("No logger")
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
