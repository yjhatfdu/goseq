package goseq

var defaultLogger *Logger

func Connect(apiEndpoint string, apiKey string) {
	defaultLogger = NewLogger(apiEndpoint, apiKey, 1024)
}

func Log(level string, messageTemplate string, params ...interface{}) {
	if defaultLogger == nil {
		return
	}
	defaultLogger.Log(level, messageTemplate, params...)
}

func Trace(messageTemplate string, params ...interface{}) {
	Log("TRACE", messageTemplate, params...)
}

func Debug(messageTemplate string, params ...interface{}) {
	Log("DEBUG", messageTemplate, params...)
}

func Info(messageTemplate string, params ...interface{}) {
	Log("INFO", messageTemplate, params...)
}

func Warn(messageTemplate string, params ...interface{}) {
	Log("WARN", messageTemplate, params...)
}

func Error(messageTemplate string, params ...interface{}) {
	Log("ERROR", messageTemplate, params...)
}

func Fatal(messageTemplate string, params ...interface{}) {
	Log("FATAL", messageTemplate, params...)
}
