package logadapter

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

// EchoLogger extend logrus.Logger
type EchoLogger struct {
	*Logger
}

// NewEchoLogger return singleton logger
func NewEchoLogger() *EchoLogger {
	return &EchoLogger{Logger: l}
}

// To logrus.Level
func toLogrusLevel(level log.Lvl) logrus.Level {
	switch level {
	case log.DEBUG:
		return logrus.DebugLevel
	case log.INFO:
		return logrus.InfoLevel
	case log.WARN:
		return logrus.WarnLevel
	case log.ERROR:
		return logrus.ErrorLevel
	}

	return logrus.InfoLevel
}

// To Echo.log.lvl
func toEchoLevel(level logrus.Level) log.Lvl {
	switch level {
	case logrus.DebugLevel:
		return log.DEBUG
	case logrus.InfoLevel:
		return log.INFO
	case logrus.WarnLevel:
		return log.WARN
	case logrus.ErrorLevel:
		return log.ERROR
	}

	return log.OFF
}

// Output return logger io.Writer
func (l *EchoLogger) Output() io.Writer {
	return l.Out
}

// SetOutput logger io.Writer
func (l *EchoLogger) SetOutput(w io.Writer) {
	l.Out = w
}

// Level return logger level
func (l *EchoLogger) Level() log.Lvl {
	return toEchoLevel(l.Logger.Level)
}

// SetLevel logger level
func (l *EchoLogger) SetLevel(v log.Lvl) {
	l.Logger.Level = toLogrusLevel(v)
}

// SetHeader logger header
// Managed by Logrus itself
// This function do nothing
func (l *EchoLogger) SetHeader(h string) {
	// do nothing
}

// Formatter return logger formatter
func (l *EchoLogger) Formatter() logrus.Formatter {
	return l.Logger.Formatter
}

// SetFormatter logger formatter
// Only support logrus formatter
func (l *EchoLogger) SetFormatter(formatter logrus.Formatter) {
	l.Logger.Formatter = formatter
}

// Prefix return logger prefix
// This function do nothing
func (l *EchoLogger) Prefix() string {
	return ""
}

// SetPrefix logger prefix
// This function do nothing
func (l *EchoLogger) SetPrefix(p string) {
	// do nothing
}

// Print output message of print level
func (l *EchoLogger) Print(i ...interface{}) {
	l.Logger.Print(i...)
}

// Printf output format message of print level
func (l *EchoLogger) Printf(format string, args ...interface{}) {
	l.Logger.Printf(format, args...)
}

// Printj output json of print level
func (l *EchoLogger) Printj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	l.Logger.Println(string(b))
}

// Debug output message of debug level
func (l *EchoLogger) Debug(i ...interface{}) {
	l.Logger.Debug(i...)
}

// Debugf output format message of debug level
func (l *EchoLogger) Debugf(format string, args ...interface{}) {
	l.Logger.Debugf(format, args...)
}

// Debugj output message of debug level
func (l *EchoLogger) Debugj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	l.Logger.Debugln(string(b))
}

// Info output message of info level
func (l *EchoLogger) Info(i ...interface{}) {
	l.Logger.Info(i...)
}

// Infof output format message of info level
func (l *EchoLogger) Infof(format string, args ...interface{}) {
	l.Logger.Infof(format, args...)
}

// Infoj output json of info level
func (l *EchoLogger) Infoj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	l.Logger.Infoln(string(b))
}

// Warn output message of warn level
func (l *EchoLogger) Warn(i ...interface{}) {
	l.Logger.Warn(i...)
}

// Warnf output format message of warn level
func (l *EchoLogger) Warnf(format string, args ...interface{}) {
	l.Logger.Warnf(format, args...)
}

// Warnj output json of warn level
func (l *EchoLogger) Warnj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	l.Logger.Warnln(string(b))
}

// Error output message of error level
func (l *EchoLogger) Error(i ...interface{}) {
	l.Logger.Error(i...)
}

// Errorf output format message of error level
func (l *EchoLogger) Errorf(format string, args ...interface{}) {
	l.Logger.Errorf(format, args...)
}

// Errorj output json of error level
func (l *EchoLogger) Errorj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	l.Logger.Errorln(string(b))
}

// Fatal output message of fatal level
func (l *EchoLogger) Fatal(i ...interface{}) {
	l.Logger.Fatal(i...)
}

// Fatalf output format message of fatal level
func (l *EchoLogger) Fatalf(format string, args ...interface{}) {
	l.Logger.Fatalf(format, args...)
}

// Fatalj output json of fatal level
func (l *EchoLogger) Fatalj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	l.Logger.Fatalln(string(b))
}

// Panic output message of panic level
func (l *EchoLogger) Panic(i ...interface{}) {
	l.Logger.Panic(i...)
}

// Panicf output format message of panic level
func (l *EchoLogger) Panicf(format string, args ...interface{}) {
	l.Logger.Panicf(format, args...)
}

// Panicj output json of panic level
func (l *EchoLogger) Panicj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	l.Logger.Panicln(string(b))
}

// NewEchoLoggerMiddleware returns a middleware that logs HTTP requests.
func NewEchoLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			// * set request_id
			id := c.Request().Header.Get(string(CorrelationIDKey))
			if id == "" {
				id = c.Request().Header.Get(string(RequestIDKey))
				if id != "" {
					c.Response().Header().Set(string(RequestIDKey), id)
				}
			} else {
				c.Response().Header().Set(string(CorrelationIDKey), id)
			}
			if id == "" {
				id = uuid.NewString()
				c.Request().Header.Set(string(CorrelationIDKey), id)
				c.Response().Header().Set(string(CorrelationIDKey), id)
			}

			// * set request_id to request context
			ctx := withCorrelationID(c.Request().Context(), id)
			request := c.Request().WithContext(ctx)
			c.SetRequest(request)

			start := time.Now()
			var err error
			var errStr string
			if err = next(c); err != nil {
				c.Error(err)
				b, _ := json.Marshal(err.Error())
				b = b[1 : len(b)-1]
				errStr = string(b)
			}
			stop := time.Now()
			reqSizeStr := req.Header.Get(echo.HeaderContentLength)
			if reqSizeStr == "" {
				reqSizeStr = "0"
			}
			reqSize, _ := strconv.ParseInt(reqSizeStr, 10, 64)

			// * log json format
			latency := stop.Sub(start)
			trace := map[string]interface{}{
				"time":       stop.Format(l.timestampFormat),
				"ip":         c.RealIP(),
				"user_agent": req.UserAgent(),
				"host":       req.Host,
				"method":     req.Method,
				"url":        req.RequestURI,
				"status":     res.Status,
				"byte_in":    reqSize,
				"byte_out":   res.Size,
				"latency":    latency.String(),
				"latency_ms": latency.Milliseconds(),
				"referer":    req.Referer(),
				"type":       LogTypeAPI,
				"request_id": getCorrelationID(ctx),
				"error":      errStr,
			}

			var buf bytes.Buffer
			b, _ := json.Marshal(trace)
			buf.Write(b)
			buf.WriteString("\n")
			if logger, ok := c.Logger().(*EchoLogger); ok {
				logger.Output().Write(buf.Bytes())
			} else {
				c.Logger().Output().Write(buf.Bytes())
			}

			return err
		}
	}
}

// LogWithEchoContext log content with echo context
// content[0] : message -> interface{},
// content[1] : log type -> string,
// content[2] : log field -> map[string]interface{}
func LogWithEchoContext(c echo.Context, content ...interface{}) {
	var logType string
	if len(content) > 1 {
		if value, ok := content[1].(string); ok && value != "" {
			logType = value
		} else {
			logType = LogTypeDebug
		}
	}

	logField := logrus.Fields{
		"type":       logType,
		"request_id": getCorrelationID(c.Request().Context()),
	}

	if len(content) > 2 {
		if maps, ok := content[2].(map[string]interface{}); ok {
			for key, value := range maps {
				logField[key] = value
			}
		}
	}

	switch logType {
	case LogTypeAPI:
		if logger, ok := c.Logger().(*EchoLogger); ok {
			logger.WithFields(logField).Info(content[0])
		} else {
			if len(content) > 2 {
				b, _ := json.Marshal(content[2])
				c.Logger().Info(content[0], ",", string(b))
			} else {
				c.Logger().Info(content[0])
			}
		}
	case LogTypeError:
		if logger, ok := c.Logger().(*EchoLogger); ok {
			logger.WithFields(logField).Error(content[0])
		} else {
			if len(content) > 2 {
				b, _ := json.Marshal(content[2])
				c.Logger().Error(content[0], ",", string(b))
			} else {
				c.Logger().Error(content[0])
			}
		}
	default:
		if logger, ok := c.Logger().(*EchoLogger); ok {
			logger.WithFields(logField).Debug(content[0])
		} else {
			if len(content) > 2 {
				b, _ := json.Marshal(content[2])
				c.Logger().Debug(content[0], ",", string(b))
			} else {
				c.Logger().Debug(content[0])
			}
		}
	}
}
