package logadapter

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

// Key is key context type
type Key string

// Exported constanst
const (
	CorrelationIDKey Key = "X-User-Correlation-Id"
	RequestIDKey     Key = "X-Request-ID"
)

// WithCorrelationID sets correlation id to context
func WithCorrelationID(parent context.Context, correlationID string) context.Context {
	return context.WithValue(parent, CorrelationIDKey, correlationID)
}

// GetCorrelationID return correlation id
func GetCorrelationID(ctx context.Context) string {
	id := ctx.Value(CorrelationIDKey)
	if id != nil {
		return id.(string)
	}
	return ""
}

// GetLogFile get file log
func GetLogFile() string {
	path := ""
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	path = dir + "/logs"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModeDir)
		if err != nil {
			panic(err)
		}
	}
	if strings.Contains(runtime.GOOS, "window") {
		path += "\\"
	} else {
		path += "/"
	}
	return path + fmt.Sprintf("log_%s.log", time.Now().Format("2006-01-02T15:04:05"))
}
