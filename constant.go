package logadapter

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

// Export HeaderKey constanst
const (
	CorrelationIDHeaderKey HeaderKey = "X-User-Correlation-Id"
	RequestIDHeaderKey     HeaderKey = "X-Request-ID"
	UserInfoHeaderKey      HeaderKey = "X-Userinfo"
)

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

// custom log format
const (
	JSONFormat LogFormat = iota
	PrettyJSONFormat
	TextFormat
)

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
