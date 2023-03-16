# go-logadapter
[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/vuduongtp/go-logadapter/blob/main/LICENSE)
![GolangVersion](https://img.shields.io/github/go-mod/go-version/vuduongtp/go-logadapter)
[![Release](https://img.shields.io/github/v/release/vuduongtp/go-logadapter)](https://github.com/vuduongtp/go-logadapter/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/vuduongtp/go-logadapter.svg)](https://pkg.go.dev/github.com/vuduongtp/go-logadapter)

<p align="center">
  <img src="https://user-images.githubusercontent.com/32934289/224491943-30e48110-ea1f-4a95-a396-5f04dd91f963.png" />
</p>

**go-logadapter provide a flexible and powerful way to handle logging in applications, and can help in debugging, monitoring, and maintaining the application's performance and behavior.**

In Go, the logging package provides a simple logging interface for sending log messages to various outputs, such as the console or a file. However, it can be useful to have more advanced logging features, such as the ability to write logs to multiple destinations, format log messages differently based on severity or context, or filter logs based on certain criteria. To accomplish these more advanced logging features, we can use **go-logadapter**. 

It's a piece of code that sits between the application and the logging package and modifies, enhances the way log messages are handled. It's customized to suit for [Echo web framework](https://github.com/labstack/echo) , [gorm](https://github.com/go-gorm/gorm) and still updating.
## Advantages of go-logadapter
- Writing logs to multiple destinations, such as a file, console.
- Formatting log messages such as JSON, pretty JSON, text.
- Filtering logs based on certain criteria, such as log level, module, type, or request_id.
- Suit specific application needs such as [echo](https://github.com/labstack/echo), [gorm](https://github.com/go-gorm/gorm)
- Can help in debugging an application by providing detailed information about the application's behavior, performance, and errors.

## Requirements

- Go 1.18+

## Getting Started

```
$ go get -u github.com/vuduongtp/go-logadapter
```
## Import
```go
import "github.com/vuduongtp/go-logadapter"
```
## Basic Example
View full example [here](https://github.com/vuduongtp/go-logadapter/blob/main/test/test.go)
### Simple example
```go
logadapter.Debug("debug message")
logadapter.Error("error message")
logadapter.Warn("warn message")
```
```
{"level":"debug","msg":"debug message","time":"2023-03-17 00:10:17.18915"}
{"level":"error","msg":"error message","time":"2023-03-17 00:10:17.18924"}
{"level":"warning","msg":"warn message","time":"2023-03-17 00:10:17.18924"}
```
### Create new logger with config
```go
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
```
```
{"level":"debug","msg":"test","time":"2023-03-17T00:14:56.763181+07:00"}
```
**Log with pertty JSON format**
```go
logadapter.SetFormatter(logadapter.PrettyJSONFormat)
logadapter.Debug("message")
```
```
{
  "level": "debug",
  "msg": "message",
  "time": "2023-03-17 00:03:07.14439"
}
```
**Log with text format**
```go
logadapter.SetFormatter(logadapter.TextFormat)
logadapter.Debug("message")
```
```
time="2023-03-17 00:03:53.74972" level=debug msg=message
```
### Set gorm logger
```go
// * set log adapter for gorm logging
db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
if err != nil {
  panic("failed to connect database")
}
db.Logger = logadapter.NewGormLogger()
```
### Set echo logger
```go
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
```
**If you really want to help us, simply Fork the project and apply for Pull Request. Thanks.**
