package api

type Logger interface {
	Debug(message string, args ...interface{})
	Notice(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warning(message string, args ...interface{})
	Error(message string, args ...interface{})
}
