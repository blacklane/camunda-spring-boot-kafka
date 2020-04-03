package middleware

import (
	"net/http"

	warsaw "github.com/blacklane/warsaw/logger"
)

func NewLoggerMiddleware(appName string) func(handler http.Handler) http.Handler {
	logger := warsaw.NewKievRequestLogger(appName)

	fn := func(next http.Handler) http.Handler {
		return logger(next.(http.HandlerFunc))
	}
	return fn
}
