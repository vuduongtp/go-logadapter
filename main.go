package logadapter

import (
	"context"
	"io"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// custom logtype
const (
	LogTypeAPI      = "api"
	LogTypeRequest  = "request"
	LogTypeResponse = "response"
	LogTypeError    = "error"
	LogTypeDebug    = "debug"
	LogTypeInfo     = "info"
	LogTypeWarn     = "warn"
	LogTypeSQL      = "sql"
	LogTypeTrace    = "trace"
)

// custom constants
const (
	DefaultTimestampFormat = "2006-01-02 15:04:05.00000"
	DefaultPrefix          = "LogAdapter_"
	DefaultSourceField     = "stack_trace"
)

// HeaderKey is key from http Header
type HeaderKey string

// Export HeaderKey constanst
const (
	CorrelationIDHeaderKey HeaderKey = "X-User-Correlation-Id"
	RequestIDHeaderKey     HeaderKey = "X-Request-ID"
	UserInfoHeaderKey      HeaderKey = "X-Userinfo"
)

// LogKey is key for log messages
type LogKey string

// Export LogKey constanst
const (
	CorrelationIDLogKey LogKey = "correlation_id"
	RequestIDLogKey     LogKey = "request_id"
	UserInfoLogKey      LogKey = "user_info"
)

// Export default LogKeyMap
var (
	DefaultLogKeys []LogKey = []LogKey{CorrelationIDLogKey, RequestIDLogKey, UserInfoLogKey}
	baseSourceDir  string
)

// LogFormat log format
type LogFormat uint32

// custom log format
const (
	JSONFormat LogFormat = iota
	PrettyJSONFormat
	TextFormat
)

// Level log level
type Level uint32

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

// Config config instance log
type Config struct {
	IsUseLogFile    bool        // set true if write to file
	FileConfig      *FileConfig // ignore if IsUseLogFile = false, set null if use default log file config
	LogLevel        Level
	LogFormat       LogFormat
	TimestampFormat string // if empty, use default timestamp format
}

// FileConfig config for write log to file
type FileConfig struct {
	Filename       string
	MaxSize        int // megabytes
	MaxBackups     int // number of log files
	MaxAge         int // days
	IsCompress     bool
	IsUseLocalTime bool
}

// Logger instance
type Logger struct {
	*log.Logger
	logFormat       LogFormat
	timestampFormat string
	logKeys         []LogKey
}

var l *Logger

func init() {
	l = New()
	sourceDir()
}

type customFormat struct {
	defaultFields map[string]interface{}
	formatter     log.Formatter
}

func (cl customFormat) Format(entry *log.Entry) ([]byte, error) {
	for k, v := range cl.defaultFields {
		entry.Data[k] = v
	}
	return cl.formatter.Format(entry)
}

// SetFormatter set logger formatter
func SetFormatter(logFormat LogFormat) { l.SetFormatter(logFormat) }

// SetFormatter set logger formatter
func (l *Logger) SetFormatter(logFormat LogFormat) {
	switch logFormat {
	case JSONFormat:
		l.Logger.SetFormatter(&log.JSONFormatter{TimestampFormat: DefaultTimestampFormat})

	case PrettyJSONFormat:
		l.Logger.SetFormatter(&log.JSONFormatter{PrettyPrint: true, TimestampFormat: DefaultTimestampFormat})

	default:
		l.Logger.SetFormatter(&log.TextFormatter{TimestampFormat: DefaultTimestampFormat})
	}

	l.logFormat = logFormat
}

// GetFormatter get logger formatter
func GetFormatter() LogFormat { return l.logFormat }

// SetTimestampFormat set timestamp format
func SetTimestampFormat(timestampFormat string) { l.SetTimestampFormat(timestampFormat) }

// SetTimestampFormat set timestamp format
func (l *Logger) SetTimestampFormat(timestampFormat string) {
	switch l.logFormat {
	case JSONFormat:
		l.Logger.SetFormatter(&log.JSONFormatter{TimestampFormat: timestampFormat})

	case PrettyJSONFormat:
		l.Logger.SetFormatter(&log.JSONFormatter{PrettyPrint: true, TimestampFormat: timestampFormat})

	default:
		l.Logger.SetFormatter(&log.TextFormatter{TimestampFormat: timestampFormat})
	}

	l.timestampFormat = timestampFormat
}

// GetTimestampFormat get timestamp format
func GetTimestampFormat() string { return l.timestampFormat }

// SetDefaultFields set default fields for all log
func SetDefaultFields(fields map[string]interface{}) { l.SetDefaultFields(fields) }

// SetDefaultFields set default fields for all log
func (l *Logger) SetDefaultFields(fields map[string]interface{}) {
	l.Logger.SetFormatter(customFormat{
		defaultFields: fields,
		formatter:     l.Formatter,
	})
}

// SetLogFile set log file, log file will be storaged in logs folder
func SetLogFile() { l.SetLogFile() }

// SetLogFileWithConfig set log file with file config
func (l *Logger) SetLogFileWithConfig(fileConfig *FileConfig) {
	if fileConfig == nil {
		fileConfig = getDefaultFileConfig()
	}

	var lumber = &lumberjack.Logger{
		Filename:   fileConfig.Filename,
		MaxSize:    fileConfig.MaxSize,
		MaxBackups: fileConfig.MaxBackups,
		MaxAge:     fileConfig.MaxAge,
		Compress:   fileConfig.IsCompress,
		LocalTime:  fileConfig.IsUseLocalTime,
	}
	writer := io.MultiWriter(os.Stdout, lumber)
	l.Logger.SetOutput(writer)
}

// SetLogFile set log file, log file will be storaged in logs folder
func (l *Logger) SetLogFile() {
	fileConfig := getDefaultFileConfig()
	l.SetLogFileWithConfig(fileConfig)
}

// SetLogConsole set log console
func SetLogConsole() { l.SetLogConsole() }

// SetLogConsole set log console
func (l *Logger) SetLogConsole() {
	l.SetOutput(os.Stdout)
}

// SetLevel set log level
func SetLevel(level Level) { l.SetLevel(level) }

// SetLevel set log level
func (l *Logger) SetLevel(level Level) {
	l.Logger.SetLevel(log.Level(level))
}

// SetLogger set logger instance
func SetLogger(logger *Logger) { l = logger }

// GetLogger set logger instance
func GetLogger() *Logger { return l }

func getDefaultFileConfig() *FileConfig {
	return &FileConfig{
		Filename:       getLogFile(),
		MaxSize:        10,
		MaxBackups:     3,
		MaxAge:         30,
		IsCompress:     false,
		IsUseLocalTime: true,
	}
}

func getDefaultConfig() *Config {
	return &Config{
		LogLevel:        DebugLevel,
		LogFormat:       JSONFormat,
		IsUseLogFile:    false,
		TimestampFormat: DefaultTimestampFormat,
	}
}

// RemoveLogKey remove a log key will not log this key
func (l *Logger) RemoveLogKey(key string) {
	for i := 0; i < len(l.logKeys); i++ {
		if strings.EqualFold(string(l.logKeys[i]), key) {
			l.logKeys = append(l.logKeys[:i], l.logKeys[i+1:]...)
		}
	}
}

// RemoveLogKey export remove a log key will not log this key
func RemoveLogKey(key string) { l.RemoveLogKey(key) }

// addLogKey add one more log key
func (l *Logger) addLogKey(key string) {
	if !logKeyExists(l.logKeys, LogKey(key)) {
		l.logKeys = append(l.logKeys, LogKey(key))
	}
}

// SetContext set log with context
func (l *Logger) SetContext(ctx context.Context) *log.Entry {
	return l.WithContext(ctx).WithFields(l.GetLogFieldFromContext(ctx))
}

// SetContext set log with context
func SetContext(ctx context.Context) *log.Entry {
	return l.SetContext(ctx)
}

// GetLogFieldFromContext gets log field from context for log field
func GetLogFieldFromContext(ctx context.Context) map[string]interface{} {
	return l.GetLogFieldFromContext(ctx)
}

// GetLogFieldFromContext gets log field from context for log field
func (l *Logger) GetLogFieldFromContext(ctx context.Context) map[string]interface{} {
	fields := make(map[string]interface{})
	for _, key := range l.logKeys {
		val := getContextKeyValue(ctx, string(key))
		if val != nil {
			fields[string(key)] = val
		}
	}

	return fields
}

// SetCustomLogField set custom log field for always log this field, return new context
func (l *Logger) SetCustomLogField(ctx context.Context, logKey string, value interface{}) context.Context {
	l.addLogKey(logKey)
	return setContextKeyValue(ctx, logKey, value)
}

// SetCustomLogField set custom log field for always log this field, return new context
func SetCustomLogField(ctx context.Context, logKey string, value interface{}) context.Context {
	return l.SetCustomLogField(ctx, logKey, value)
}

// NewWithConfig returns a logger instance with custom configuration
func NewWithConfig(config *Config) *Logger {
	if config == nil {
		config = getDefaultConfig()
	}
	logger := log.New()
	l := &Logger{Logger: logger}
	l.logFormat = config.LogFormat
	l.SetFormatter(config.LogFormat)
	if len(config.TimestampFormat) > 0 {
		l.SetTimestampFormat(config.TimestampFormat)
	}
	if config.IsUseLogFile == true {
		l.SetLogFileWithConfig(config.FileConfig)
	} else {
		l.SetLogConsole()
	}
	l.SetLevel(config.LogLevel)
	l.logKeys = DefaultLogKeys

	return l
}

// New returns a logger instance with default configuration
func New() *Logger {
	config := getDefaultConfig()
	logger := log.New()
	l := &Logger{Logger: logger}
	l.SetFormatter(config.LogFormat)
	l.logFormat = config.LogFormat
	if len(config.TimestampFormat) > 0 {
		l.SetTimestampFormat(config.TimestampFormat)
	}
	if config.IsUseLogFile == true {
		l.SetLogFileWithConfig(config.FileConfig)
	} else {
		l.SetLogConsole()
	}
	l.SetLevel(config.LogLevel)
	l.logKeys = DefaultLogKeys

	return l
}

// Trace log with trace level
func Trace(args ...interface{}) {
	l.Trace(args...)
}

// Debug log with debug level
func Debug(args ...interface{}) {
	l.Debug(args...)
}

// Info log with info level
func Info(args ...interface{}) {
	l.Info(args...)
}

// Warn log with warn level
func Warn(args ...interface{}) {
	field := log.Fields{DefaultSourceField: getCaller()}
	l.WithFields(field).Warn(args...)
}

// Error log with error level
func Error(args ...interface{}) {
	field := log.Fields{DefaultSourceField: getCaller()}
	l.WithFields(field).Error(args...)
}

// Fatal log with fatal level
func Fatal(args ...interface{}) {
	field := log.Fields{DefaultSourceField: getCaller()}
	l.WithFields(field).Fatal(args...)
}

// Panic log with panic level
func Panic(args ...interface{}) {
	field := log.Fields{DefaultSourceField: getCaller()}
	l.WithFields(field).Panic(args...)
}

// TraceWithContext log with trace level
func TraceWithContext(ctx context.Context, args ...interface{}) {
	l.SetContext(ctx).Trace(args...)
}

// DebugWithContext log with debug level
func DebugWithContext(ctx context.Context, args ...interface{}) {
	l.SetContext(ctx).Debug(args...)
}

// InfoWithContext log with info level
func InfoWithContext(ctx context.Context, args ...interface{}) {
	l.SetContext(ctx).Info(args...)
}

// WarnWithContext log with warn level
func WarnWithContext(ctx context.Context, args ...interface{}) {
	field := log.Fields{DefaultSourceField: getCaller()}
	l.SetContext(ctx).WithFields(field).Warn(args...)
}

// ErrorWithContext log with error level
func ErrorWithContext(ctx context.Context, args ...interface{}) {
	field := log.Fields{DefaultSourceField: getCaller()}
	l.SetContext(ctx).WithFields(field).Error(args...)
}

// FatalWithContext log with fatal level
func FatalWithContext(ctx context.Context, args ...interface{}) {
	field := log.Fields{DefaultSourceField: getCaller()}
	l.SetContext(ctx).WithFields(field).Fatal(args...)
}

// PanicWithContext log with panic level
func PanicWithContext(ctx context.Context, args ...interface{}) {
	field := log.Fields{DefaultSourceField: getCaller()}
	l.SetContext(ctx).WithFields(field).Panic(args...)
}

// LogWithContext log content with context
// content[0] : message -> interface{},
// content[1] : log type -> string,
// content[2] : log field -> map[string]interface{}
func LogWithContext(ctx context.Context, content ...interface{}) {
	var logType string
	if len(content) > 1 {
		if value, ok := content[1].(string); ok && value != "" {
			logType = value
		} else {
			logType = LogTypeDebug
		}
	}

	logFields := mergeLogFields(GetLogFieldFromContext(ctx), map[string]interface{}{"type": logType})

	if len(content) > 2 {
		if maps, ok := content[2].(map[string]interface{}); ok {
			logFields = mergeLogFields(logFields, maps)
		}
	}

	switch logType {
	case LogTypeAPI:
		l.Logger.WithFields(logFields).Info(content[0])
	case LogTypeError:
		field := log.Fields{DefaultSourceField: getCaller()}
		l.Logger.WithFields(mergeLogFields(logFields, field)).Error(content[0])
	case LogTypeInfo:
		l.Logger.WithFields(logFields).Info(content[0])
	case LogTypeWarn:
		field := log.Fields{DefaultSourceField: getCaller()}
		l.Logger.WithFields(mergeLogFields(logFields, field)).Warn(content[0])
	case LogTypeRequest, LogTypeResponse:
		l.Logger.WithFields(logFields).Info(content[0])
	default:
		l.Logger.WithFields(logFields).Debug(content[0])
	}
}
