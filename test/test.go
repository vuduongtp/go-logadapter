package main

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/vuduongtp/go-logadapter"
	"gorm.io/gorm"
)

func main() {
	testCreateInstance()
}

func testCreateInstance() {
	config := &logadapter.Config{
		LogLevel:     logadapter.DebugLevel,
		LogFormat:    logadapter.JSONFormat,
		IsUseLogFile: true,
	}
	logger := logadapter.NewLoggerWithConfig(config)
	logger.Debug("test")
}

func testGormAdapter() {
	config := &logadapter.Config{
		LogLevel:     logadapter.DebugLevel,
		LogFormat:    logadapter.JSONFormat,
		IsUseLogFile: true,
	}
	logger := logadapter.NewLoggerWithConfig(config)

	// * set log adapter for gorm logging
	gormConfig := new(gorm.Config)
	gormConfig.PrepareStmt = true
	gormLogAdapter := logadapter.NewGormLogAdapter(logger)
	gormLogAdapter.SlowThreshold = time.Second
	gormLogAdapter.SourceField = logadapter.DefaultGormSourceField
	gormConfig.Logger = gormLogAdapter

}

func testEchoAdapter() {
	config := &logadapter.Config{
		LogLevel:     logadapter.DebugLevel,
		LogFormat:    logadapter.JSONFormat,
		IsUseLogFile: true,
	}
	logger := logadapter.NewLoggerWithConfig(config)

	e := echo.New()
	// * set log adapter for echo instance
	e.Logger = logadapter.NewEchoLogAdapter(logger)

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
