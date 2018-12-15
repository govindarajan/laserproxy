package logger

import "log"

type LogLevel byte

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	CRITICAL
)

func (lev LogLevel) String() string {
	switch lev {
	case CRITICAL:
		return "CRITICAL"
	case ERROR:
		return "ERROR"
	case WARN:
		return "WARN"
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}

func Log(lev LogLevel, msg string) {
	log.Println(lev.String(), msg)
}

func LogDebug(msg string) {
	Log(DEBUG, msg)
}

func LogInfo(msg string) {
	Log(INFO, msg)
}

func LogWarn(msg string) {
	Log(WARN, msg)
}

func LogError(msg string) {
	Log(ERROR, msg)
}

func LogCritical(msg string) {
	Log(CRITICAL, msg)
}
