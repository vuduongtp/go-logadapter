package logadapter

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
)

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

func sourceDir() {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)

	s := filepath.Dir(dir)
	if filepath.Base(s) != "go-logadapter" {
		s = dir
	}
	baseSourceDir = filepath.ToSlash(s) + "/"
}

// setContextKeyValue sets key value to context
func setContextKeyValue(parent context.Context, key, value interface{}) context.Context {
	return context.WithValue(parent, fmt.Sprintf("%s%s", DefaultPrefix, key), value)
}

// getContextKeyValue gets key value from context
func getContextKeyValue(ctx context.Context, key interface{}) interface{} {
	return ctx.Value(fmt.Sprintf("%s%s", DefaultPrefix, key))
}

// generateCorrelationID generate correlation ID by snowflake and return string
func generateCorrelationID() string {
	// Create a new Node with a Node number of 1
	id := uuid.NewString()

	return strings.ReplaceAll(id, "-", "")
}

// mergeLogFields merge many map to 1
func mergeLogFields(maps ...map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})

	for _, m := range maps {
		for k, v := range m {
			merged[k] = v
		}
	}

	return merged
}

func logKeyExists(arr []LogKey, key LogKey) bool {
	for _, s := range arr {
		if s == key {
			return true
		}
	}
	return false
}

func getCaller() string {
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && (!strings.HasPrefix(file, baseSourceDir) || strings.HasSuffix(file, "_test.go")) {
			return fmt.Sprintf("%s:%d", file, line)
		}
	}

	return ""
}
