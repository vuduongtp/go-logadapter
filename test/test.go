package main

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/vuduongtp/go-logadapter"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	testSimpleLog()
	testNewWithConfig()
	testSetPrettyJSONFormat()
	testSetTextFormat()
	testSetTimestampFormat()
	testSetLogConsole()
	testSetLogFile()
}

func testSimpleLog() {
	logadapter.Debug("debug message")
	logadapter.Error("error message")
	logadapter.Warn("warn message")
}

func testNewWithConfig() {
	logadapter.SetLogger(logadapter.NewWithConfig(&logadapter.Config{
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
	}))
	logadapter.Debug("test")
}

func testSetPrettyJSONFormat() {
	logadapter.SetFormatter(logadapter.PrettyJSONFormat)
	logadapter.Debug("message")
}

func testSetTextFormat() {
	logadapter.SetFormatter(logadapter.TextFormat)
	logadapter.Debug("message")
}

func testSetTimestampFormat() {
	logadapter.SetTimestampFormat(time.RFC1123)
	logadapter.Debug("message")
}

func testSetLogConsole() {
	logadapter.SetLogConsole()
	logadapter.Debug("message")
}

func testSetLogFile() {
	logadapter.SetLogFile(&logadapter.FileConfig{
		Filename:       "logs",
		MaxSize:        50,
		MaxBackups:     10,
		MaxAge:         30,
		IsCompress:     false,
		IsUseLocalTime: false,
	})
	logadapter.Debug("message")
}

func testGormAdapter() {
	// * set log adapter for gorm logging
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.Logger = logadapter.NewGormLogger()
}

func testEchoAdapter() {
	e := echo.New()
	// * set log adapter for echo instance
	e.Logger = logadapter.NewEchoLogger()

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
