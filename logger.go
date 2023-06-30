package redisson

type Logger interface {
	Debug(message string, args ...interface{})
	Notice(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warning(message string, args ...interface{})
	Error(message string, args ...interface{})
}

// region - logger

func (r *redis) Debug(message string, args ...interface{}) {
	if r.logger != nil {
		r.logger.Debug(message, args...)
	}
}
func (r *redis) Notice(message string, args ...interface{}) {
	if r.logger != nil {
		r.logger.Notice(message, args...)
	}
}
func (r *redis) Info(message string, args ...interface{}) {
	if r.logger != nil {
		r.logger.Info(message, args...)
	}
}
func (r *redis) Warning(message string, args ...interface{}) {
	if r.logger != nil {
		r.logger.Warning(message, args...)
	}
}
func (r *redis) Error(message string, args ...interface{}) {
	if r.logger != nil {
		r.logger.Error(message, args...)
	}
}

// endregion
