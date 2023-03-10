package logadapter

import (
	"io"
	"os"

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
	LogTypeSQL      = "sql"
	LogTypeTrace    = "trace"
)

// custom constants
const (
	DefaultTimestampFormat = "2006-01-02 15:04:05.00000"
	DefaultGormSourceField = "source"
)

// Key is key context type
type Key string

// Exported constanst
const (
	CorrelationIDKey Key = "X-User-Correlation-Id"
	RequestIDKey     Key = "X-Request-ID"
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
	logFormat LogFormat
}

// SetFormatter logger formatter
func (l *Logger) SetFormatter(logFormat LogFormat) {
	l.logFormat = logFormat
	switch logFormat {
	case JSONFormat:
		l.Logger.SetFormatter(&log.JSONFormatter{TimestampFormat: DefaultTimestampFormat})

	case PrettyJSONFormat:
		l.Logger.SetFormatter(&log.JSONFormatter{PrettyPrint: true, TimestampFormat: DefaultTimestampFormat})

	default:
		l.Logger.SetFormatter(&log.TextFormatter{TimestampFormat: DefaultTimestampFormat})
	}
}

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
}

// SetLogFile set log file
func (l *Logger) SetLogFile(fileConfig *FileConfig) {
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

// SetLogConsole set log console
func (l *Logger) SetLogConsole() {
	l.SetOutput(os.Stdout)
}

// SetLevel set log level
func (l *Logger) SetLevel(level Level) {
	l.Logger.SetLevel(log.Level(level))
}

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

// NewLoggerWithConfig returns a logger instance with custom configuration
func NewLoggerWithConfig(config *Config) *Logger {
	if config == nil {
		config = getDefaultConfig()
	}
	logger := log.New()
	loggerInstance := &Logger{Logger: logger}
	loggerInstance.logFormat = config.LogFormat
	loggerInstance.SetFormatter(config.LogFormat)
	if len(config.TimestampFormat) > 0 {
		loggerInstance.SetTimestampFormat(config.TimestampFormat)
	}
	if config.IsUseLogFile == true {
		loggerInstance.SetLogFile(config.FileConfig)
	} else {
		loggerInstance.SetLogConsole()
	}
	loggerInstance.SetLevel(config.LogLevel)
	loggerInstance.Info("Logger instance has been successfully initialized")

	return loggerInstance
}

// NewLogger returns a logger instance with default configuration
func NewLogger() *Logger {
	config := getDefaultConfig()
	logger := log.New()
	loggerInstance := &Logger{Logger: logger}
	loggerInstance.SetFormatter(config.LogFormat)
	loggerInstance.logFormat = config.LogFormat
	if len(config.TimestampFormat) > 0 {
		loggerInstance.SetTimestampFormat(config.TimestampFormat)
	}
	if config.IsUseLogFile == true {
		loggerInstance.SetLogFile(config.FileConfig)
	} else {
		loggerInstance.SetLogConsole()
	}
	loggerInstance.SetLevel(config.LogLevel)
	loggerInstance.Info("Logger instance has been successfully initialized")

	return loggerInstance
}
