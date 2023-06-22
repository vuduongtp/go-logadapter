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
- Filtering logs based on certain criteria, such as log level, module, type, or request_id, correlation_id.
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
{"level":"debug","msg":"debug message","time":"2023-06-22 21:27:08.97942"}
{"level":"error","msg":"error message","source":"/go-logadpater:25","time":"2023-06-22 21:27:08.97951"}
{"level":"warning","msg":"warn message","source":"/go-logadpater:25","time":"2023-06-22 21:27:08.97952"}
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
**Add custome log field**
```go
ctx := context.Background()
ctx = logadapter.SetCustomLogField(ctx, "test", "test")
logadapter.InfoWithContext(ctx, "This is message 1")

logadapter.RemoveLogKey("test") // remove this key from log messages
logadapter.InfoWithContext(ctx, "This is message 2")
```
```
{"level":"info","test":"test","msg":"This is message 1","time":"2023-06-21 17:18:14.49578"}
{"level":"info","msg":"This is message 2","time":"2023-06-21 17:18:14.49586"}
```
### Set gorm logger
```go
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
```
```
{"level":"debug","database_name":"test","msg":"","time":"2023-06-21 16:53:53.14278","row":1,"query":"SELECT * FROM users ORDER BY id LIMIT 1","latency":"63.995042ms","latency_ms":63,"type":"sql"}
```
### Set echo logger
```go
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
```
```
   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v4.10.2
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                  O\
⇨ http server started on [::]:1323

-> go to http://localhost:1323/

{"correlation_id":"181e60c9d7b144a7a3960852b17efa45","level":"info","msg":"Message: Hello, World!","request_id":"ef4a720b-8af2-45b0-bf0b-4bdcb2424bd9","time":"2023-06-21 17:03:01.92469"}
{"byte_in":0,"byte_out":13,"correlation_id":"181e60c9d7b144a7a3960852b17efa45","host":"localhost:1323","ip":"127.0.0.1","latency":"427.875µs","latency_ms":0,"level":"info","method":"GET","msg":"","referer":"","request_id":"ef4a720b-8af2-45b0-bf0b-4bdcb2424bd9","status":200,"time":"2023-06-21 17:03:01.92513","type":"api","uri":"/","url":"/","user_agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/112.0"}
```
**If you really want to help us, simply Fork the project and apply for Pull Request. Thanks.**
