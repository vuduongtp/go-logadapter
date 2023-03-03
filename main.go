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

// LogFormat log format
type LogFormat uint32

// custom log format
const (
	JSONFormat LogFormat = iota
	JSONFormatIndent
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
	IsUseLogFile bool        // set true if write to file
	FileConfig   *FileConfig // ignore if IsUseLogFile = false, set null if use default log file config
	LogLevel     Level
	LogFormat    LogFormat
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

func getDefaultFileConfig() *FileConfig {
	return &FileConfig{
		Filename:       GetLogFile(),
		MaxSize:        10,
		MaxBackups:     3,
		MaxAge:         30,
		IsCompress:     false,
		IsUseLocalTime: true,
	}
}

// NewLoggerInstance returns a logrus instance
func NewLoggerInstance(config Config) *log.Logger {
	logger := log.New()

	switch config.LogFormat {
	case JSONFormat:
		logger.SetFormatter(&log.JSONFormatter{})

	case JSONFormatIndent:
		logger.SetFormatter(&log.JSONFormatter{PrettyPrint: true})

	default:
		logger.SetFormatter(&log.TextFormatter{})
	}

	if config.IsUseLogFile == true {
		if config.FileConfig == nil {
			config.FileConfig = getDefaultFileConfig()
		}

		var lumber = &lumberjack.Logger{
			Filename:   config.FileConfig.Filename,
			MaxSize:    config.FileConfig.MaxSize,
			MaxBackups: config.FileConfig.MaxBackups,
			MaxAge:     config.FileConfig.MaxAge,
			Compress:   config.FileConfig.IsCompress,
			LocalTime:  config.FileConfig.IsUseLocalTime,
		}
		writer := io.MultiWriter(os.Stdout, lumber)
		logger.SetOutput(writer)
	} else {
		logger.SetOutput(os.Stdout)
	}

	logger.SetLevel(log.Level(config.LogLevel))

	logger.Info("Logger initialization successful")

	return logger
}
