package logger

type Logger interface {
	Debug(msg string)
	Debugf(msg string, fields map[string]interface{})
	Info(msg string)
	Warn(msg string)
	Error(err error, msg string)
	Fatal(err error, msg string)
	With(fields map[string]interface{}) Logger
	Component(name string) Logger
}
