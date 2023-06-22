package main

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/vuduongtp/go-logadapter"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	testSetCustomLogField()
	testSimpleLog()
	testNewWithConfig()
	testSetPrettyJSONFormat()
	testSetTextFormat()
	testSetTimestampFormat()
	testSetLogConsole()
	testSetLogFile()
	testEchoAdapter()
}

func testSimpleLog() {
	logadapter.Debug("debug message")
	logadapter.Error("error message")
	logadapter.Warn("warn message")

	// {"level":"debug","msg":"debug message","time":"2023-06-22 21:40:55.54591"}
	// {"level":"error","msg":"error message","source":"/go-logadpater:25","time":"2023-06-22 21:40:55.54600"}
	// {"level":"warning","msg":"warn message","source":"/go-logadpater:25","time":"2023-06-22 21:40:55.54601"}
}

func testSetCustomLogField() {
	ctx := context.Background()
	ctx = logadapter.SetCustomLogField(ctx, "test", "test")
	logadapter.InfoWithContext(ctx, "This is message 1")

	logadapter.RemoveLogKey("test") // remove this key from log messages
	logadapter.InfoWithContext(ctx, "This is message 2")

	// {"level":"info","msg":"This is message 1","test":"test","time":"2023-06-22 21:41:35.48296"}
	// {"level":"info","msg":"This is message 2","time":"2023-06-22 21:41:35.48305"}
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

	// {"level":"debug","msg":"test","time":"2023-06-22T21:41:58.749418+07:00"}
}

func testSetPrettyJSONFormat() {
	logadapter.SetFormatter(logadapter.PrettyJSONFormat)
	logadapter.Debug("message")

	// {
	// 	"level": "debug",
	// 	"msg": "message",
	// 	"time": "2023-06-22 21:42:21.73744"
	//   }
}

func testSetTextFormat() {
	logadapter.SetFormatter(logadapter.TextFormat)
	logadapter.Debug("message")

	// time="2023-06-22 21:42:44.43361" level=debug msg=message
}

func testSetTimestampFormat() {
	logadapter.SetTimestampFormat(time.RFC1123)
	logadapter.Debug("message")

	// {"level":"debug","msg":"message","time":"Thu, 22 Jun 2023 21:43:01 +07"}
}

func testSetLogConsole() {
	logadapter.SetLogConsole()
	logadapter.Debug("message")

	// {"level":"debug","msg":"message","time":"2023-06-22 21:43:20.33647"}
}

func testSetLogFile() {
	logadapter.SetLogFile()
	logadapter.Debug("message")
}

func testGormAdapter() {
	isDebug := true
	// set log adapter for gorm logging
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if isDebug {
		db.Logger = logadapter.NewGormLogger().LogMode(logger.Info)
	} else {
		db.Logger = logadapter.NewGormLogger().LogMode(logger.Silent)
	}

	// set one more log key
	ctx := context.Background()
	ctx = logadapter.SetCustomLogField(ctx, "database_name", "test")

	type User interface{}
	user := new(User)
	db.WithContext(ctx).First(user)

	// {"level":"debug","database_name":"test","msg":"","time":"2023-06-21 16:53:53.14278","row":1,"query":"SELECT * FROM users ORDER BY id LIMIT 1","latency":"63.995042ms","latency_ms":63,"type":"sql"}
}

func testEchoAdapter() {
	isDebug := true
	e := echo.New()
	// set log adapter for echo instance
	e.Logger = logadapter.NewEchoLogger()
	if isDebug {
		e.Logger.SetLevel(log.DEBUG)
	} else {
		e.Logger.SetLevel(log.ERROR)
	}

	// use log adapter middleware for echo web framework
	e.Use(logadapter.NewEchoLoggerMiddleware())

	e.GET("/", func(c echo.Context) error {
		logadapter.InfoWithContext(c.Request().Context(), "Message: ", "Hello, World!")
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))

	// 	  ____    __
	// 	 / __/___/ /  ___
	//  / _// __/ _ \/ _ \
	// /___/\__/_//_/\___/ v4.10.2
	// High performance, minimalist Go web framework
	// https://echo.labstack.com
	// ____________________________________O/_______
	// 									O\
	// ⇨ http server started on [::]:1323
	//
	// -> go to http://localhost:1323/
	//
	// {"correlation_id":"181e60c9d7b144a7a3960852b17efa45","level":"info","msg":"Message: Hello, World!","request_id":"ef4a720b-8af2-45b0-bf0b-4bdcb2424bd9","time":"2023-06-21 17:03:01.92469"}
	// {"byte_in":0,"byte_out":13,"correlation_id":"181e60c9d7b144a7a3960852b17efa45","host":"localhost:1323","ip":"127.0.0.1","latency":"427.875µs","latency_ms":0,"level":"info","method":"GET","msg":"","referer":"","request_id":"ef4a720b-8af2-45b0-bf0b-4bdcb2424bd9","status":200,"time":"2023-06-21 17:03:01.92513","type":"api","uri":"/","url":"/","user_agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/112.0"}
}
