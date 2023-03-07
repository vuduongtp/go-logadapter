package logadapter

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

// withCorrelationID sets correlation id to context
func withCorrelationID(parent context.Context, correlationID string) context.Context {
	return context.WithValue(parent, CorrelationIDKey, correlationID)
}

// getCorrelationID return correlation id
func getCorrelationID(ctx context.Context) string {
	id := ctx.Value(CorrelationIDKey)
	if id != nil {
		return id.(string)
	}
	return ""
}

// getLogFile get file log
func getLogFile() string {
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
