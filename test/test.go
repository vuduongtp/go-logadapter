package main

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/vuduongtp/go-logadapter"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	testCreateInstance()
}

func testCreateInstance() {
	logger := logadapter.NewLoggerWithConfig(&logadapter.Config{
		LogLevel:        logadapter.DebugLevel,
		LogFormat:       logadapter.JSONFormat,
		TimestampFormat: time.RFC3339Nano,
		IsUseLogFile:    true,
		FileConfig: &logadapter.FileConfig{
			Filename:       "logs",
			MaxSize:        50,
			MaxBackups:     10,
			MaxAge:         30,
			IsCompress:     false,
			IsUseLocalTime: false,
		},
	})
	logger.Debug("test")
}

func testGormAdapter() {
	// * set log adapter for gorm logging
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.Logger = logadapter.NewGormLogAdapter(logadapter.NewLogger())
}

func testEchoAdapter() {
	e := echo.New()
	// * set log adapter for echo instance
	e.Logger = logadapter.NewEchoLogAdapter(logadapter.NewLogger())

	// * use log adapter middleware for echo web framework
	e.Use(logadapter.NewEchoLoggerMiddleware())

	// * log with echo context for log request_id
	echoContext := e.AcquireContext() // example echo context, should be replaced with echo.Request().Context()
	logadapter.LogWithEchoContext(echoContext, "this is message", logadapter.LogTypeDebug, map[string]interface{}{
		"field_name": "this is log field",
	}) // log message with extend field
	logadapter.LogWithEchoContext(echoContext, "this is message 2", logadapter.LogTypeError) // log message error
	logadapter.LogWithEchoContext(echoContext, "this is message 3")                          // log message debug
}
